package database

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrAssessmentConflict = errors.New("assessment is not open for this action")

const behavioralOnly = `jsonb_array_length(a.configuration_snapshot->'modules') > 0
	AND NOT EXISTS (
		SELECT 1 FROM jsonb_array_elements(a.configuration_snapshot->'modules') m
		WHERE m #>> '{configuration,definition,type}' IS DISTINCT FROM 'behavioral'
		  AND (m #> '{configuration,definition,behavioral}' = 'true'::jsonb) IS NOT TRUE
	)`

const behavioralQuestion = `coalesce(
	a.configuration_snapshot #>> ARRAY['modules', (q.module_ordinal - 1)::text, 'configuration', 'definition', 'type'] = 'behavioral'
	OR a.configuration_snapshot #> ARRAY['modules', (q.module_ordinal - 1)::text, 'configuration', 'definition', 'behavioral'] = 'true'::jsonb,
	false)`

// Snapshot ordinals are stable for existing attempts, including selected library questions.
const snapshotQuestions = `
	SELECT md5(a.id::text || ':' || m.ordinal || ':' || q.ordinal)::uuid AS id,
		m.item->>'module_id' AS module_id, q.item AS question,
		m.ordinal AS module_ordinal, q.ordinal AS question_ordinal
	FROM jsonb_array_elements(a.selected_content) WITH ORDINALITY m(item, ordinal)
	CROSS JOIN LATERAL jsonb_array_elements(m.item->'questions') WITH ORDINALITY q(item, ordinal)
`

type AssessmentRepository struct{ pool *pgxpool.Pool }

func NewAssessmentRepository(pool *pgxpool.Pool) *AssessmentRepository {
	return &AssessmentRepository{pool: pool}
}

type AssessmentAttempt struct {
	ID               string          `json:"id"`
	InstanceID       string          `json:"instanceId"`
	StageExecutionID string          `json:"stageExecutionId"`
	Status           string          `json:"status"`
	StartedAt        *time.Time      `json:"startedAt"`
	Deadline         time.Time       `json:"deadline"`
	SubmittedAt      *time.Time      `json:"submittedAt"`
	AutoSubmitted    bool            `json:"autoSubmitted"`
	Questions        json.RawMessage `json:"questions"`
}

func (r *AssessmentRepository) CandidateAttempt(ctx context.Context, candidate, id string, start, submit bool) (AssessmentAttempt, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return AssessmentAttempt{}, err
	}
	defer tx.Rollback(ctx)
	var attemptID string
	err = tx.QueryRow(ctx, `
		SELECT a.id::text FROM assessment_attempts a
		JOIN stage_executions se ON se.id = a.stage_execution_id
		JOIN pipeline_executions pe ON pe.id = se.pipeline_execution_id
		WHERE a.candidate_id = $1::uuid AND
			(CASE WHEN $3 THEN a.instance_id ELSE a.id END) = $2::uuid
			AND (NOT ($3 OR $4) OR pe.status = 'active')
			AND (NOT ($3 OR $4) OR a.submitted_at IS NOT NULL
				OR se.status IN ('available', 'in_progress'))
		FOR UPDATE OF a`, candidate, id, start, submit).Scan(&attemptID)
	if err != nil {
		return AssessmentAttempt{}, err
	}
	if start || submit {
		if _, err = tx.Exec(ctx, `
			UPDATE assessment_attempts a
			SET status = CASE
				WHEN a.deadline <= statement_timestamp() AND a.started_at IS NULL THEN 'missed'
				WHEN $2 AND a.status = 'available' AND a.deadline > statement_timestamp() THEN 'in_progress'
				WHEN a.started_at IS NOT NULL AND ($3 OR a.deadline <= statement_timestamp())
					THEN CASE WHEN `+behavioralOnly+` THEN 'submitted' ELSE 'pending_grading' END
				ELSE a.status END,
				started_at = CASE WHEN $2 AND a.status = 'available' AND a.deadline > statement_timestamp()
					THEN statement_timestamp() ELSE a.started_at END,
				deadline = CASE WHEN $2 AND a.status = 'available' AND a.deadline > statement_timestamp()
					THEN least(a.deadline, statement_timestamp() + (
						SELECT min((m->>'active_time_limit')::interval)
						FROM jsonb_array_elements(a.configuration_snapshot->'modules') m
					)) ELSE a.deadline END,
				submitted_at = CASE WHEN a.started_at IS NOT NULL AND ($3 OR a.deadline <= statement_timestamp())
					THEN least(a.deadline, statement_timestamp()) END,
				auto_submitted = a.started_at IS NOT NULL AND a.deadline <= statement_timestamp()
			WHERE a.id = $1::uuid AND a.status IN ('available', 'in_progress')
		`, attemptID, start, submit); err != nil {
			return AssessmentAttempt{}, err
		}
		if _, err = tx.Exec(ctx, `
			UPDATE assessment_instances i SET status = a.status, started_at = a.started_at AT TIME ZONE 'UTC',
				deadline = a.deadline AT TIME ZONE 'UTC'
			FROM assessment_attempts a WHERE a.id = $1::uuid AND i.id = a.instance_id
		`, attemptID); err != nil {
			return AssessmentAttempt{}, err
		}
	}
	var attempt AssessmentAttempt
	err = tx.QueryRow(ctx, `
		SELECT a.id::text, a.instance_id::text, a.stage_execution_id::text, a.status,
			a.started_at, a.deadline, a.submitted_at, a.auto_submitted,
			CASE WHEN a.started_at IS NULL THEN '[]'::jsonb ELSE coalesce((
				SELECT jsonb_agg(jsonb_strip_nulls(jsonb_build_object(
					'id', q.id, 'prompt', q.question->'prompt', 'responseFormat', q.question->'format',
					'options', q.question->'options', 'position', q.question_ordinal - 1))
					ORDER BY q.module_ordinal, q.question_ordinal)
				FROM (`+snapshotQuestions+`) q
			), '[]'::jsonb) END
		FROM assessment_attempts a WHERE a.id = $1::uuid`, attemptID).
		Scan(&attempt.ID, &attempt.InstanceID, &attempt.StageExecutionID, &attempt.Status,
			&attempt.StartedAt, &attempt.Deadline, &attempt.SubmittedAt, &attempt.AutoSubmitted, &attempt.Questions)
	if err != nil {
		return AssessmentAttempt{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return AssessmentAttempt{}, err
	}
	if (start && attempt.Status != "in_progress") || (submit && attempt.SubmittedAt == nil) {
		return AssessmentAttempt{}, ErrAssessmentConflict
	}
	return attempt, nil
}

func (r *AssessmentRepository) SaveAnswer(ctx context.Context, candidate, attempt, question string, response json.RawMessage) (time.Time, error) {
	var savedAt time.Time
	err := r.pool.QueryRow(ctx, `
		UPDATE assessment_attempts a
		SET saved_responses = jsonb_set(a.saved_responses, ARRAY[$3::uuid::text], $4::jsonb)
		WHERE a.id = $2::uuid AND a.candidate_id = $1::uuid
			AND a.status = 'in_progress' AND a.deadline > clock_timestamp()
			AND EXISTS (
				SELECT 1 FROM stage_executions se JOIN pipeline_executions pe ON pe.id = se.pipeline_execution_id
				WHERE se.id = a.stage_execution_id AND pe.status = 'active'
					AND se.status IN ('available', 'in_progress')
			)
			AND EXISTS (
				SELECT 1 FROM (`+snapshotQuestions+`) q WHERE q.id = $3::uuid
				AND ($4::jsonb = 'null'::jsonb OR CASE q.question->>'format'
					WHEN 'short_text' THEN jsonb_typeof($4::jsonb) = 'string'
					WHEN 'long_text' THEN jsonb_typeof($4::jsonb) = 'string'
					WHEN 'code' THEN jsonb_typeof($4::jsonb->'code') = 'string'
						AND jsonb_typeof($4::jsonb->'language') = 'string'
						AND nullif(btrim($4::jsonb->>'language'), '') IS NOT NULL
					WHEN 'video' THEN jsonb_typeof($4::jsonb->'videoUrl') = 'string'
						AND $4::jsonb->>'videoUrl' ~ '^https?://[^[:space:]]+$'
					WHEN 'multiple_choice' THEN CASE WHEN jsonb_typeof($4::jsonb) = 'array' THEN (
						SELECT count(*) = count(DISTINCT selection) AND coalesce(bool_and(
							CASE WHEN jsonb_typeof(selection) = 'number' THEN
								selection::numeric = trunc(selection::numeric)
								AND selection::numeric >= 0
								AND selection::numeric < jsonb_array_length(q.question->'options')
							ELSE false END), true)
						FROM jsonb_array_elements($4::jsonb) selection
					) ELSE false END
					ELSE false END)
			)
		RETURNING clock_timestamp()
	`, candidate, attempt, question, response).Scan(&savedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, ErrAssessmentConflict
	}
	return savedAt, err
}

// QuestionEvidence is keyed to the immutable question projection for this attempt.
// Evaluator integrations call FinalizeEvidence directly. Human callers are scoped
// to their posting and can only provide evidence for human-evaluated questions.
type QuestionEvidence struct {
	QuestionID  string `json:"questionId"`
	PassedCases *int   `json:"passedCases,omitempty"`
	TotalCases  *int   `json:"totalCases,omitempty"`
	Score       *int   `json:"score,omitempty"`
	Notes       string `json:"notes,omitempty"`
}

func (r *AssessmentRepository) FinalizeEvidence(ctx context.Context, attemptID, company string, evidence []QuestionEvidence) error {
	if len(evidence) == 0 {
		return ErrAssessmentConflict
	}
	for _, item := range evidence {
		if item.Score != nil && (*item.Score < 0 || *item.Score > 100) {
			return ErrAssessmentConflict
		}
	}
	evidence = slices.Clone(evidence)
	slices.SortFunc(evidence, func(a, b QuestionEvidence) int {
		return strings.Compare(a.QuestionID, b.QuestionID)
	})
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var lockedID string
	if err := tx.QueryRow(ctx, `SELECT id::text FROM assessment_attempts WHERE id = $1::uuid FOR UPDATE`,
		attemptID).Scan(&lockedID); err != nil {
		return err
	}
	var finalized bool
	err = tx.QueryRow(ctx, `
		WITH finalized AS (
			UPDATE assessment_attempts a
			SET status = 'scored', final_evidence = $2::jsonb, evaluated_at = statement_timestamp(),
				final_score = (
					SELECT CASE WHEN count(*) = count(item->>'score')
						THEN round(avg((item->>'score')::numeric))::int END
					FROM jsonb_array_elements($2::jsonb) item
				)
			WHERE a.id = $1::uuid AND a.status IN ('pending_grading', 'submitted')
				AND a.submitted_at IS NOT NULL
				AND EXISTS (
					SELECT 1 FROM stage_executions se
					JOIN pipeline_executions pe ON pe.id = se.pipeline_execution_id
					JOIN applications app ON app.id = pe.application_id
					JOIN job_postings j ON j.id = app.posting_id
					WHERE se.id = a.stage_execution_id AND pe.status = 'active'
						AND ($3 = '' OR j.company_id::text = $3)
				)
				AND jsonb_array_length($2::jsonb) = (SELECT count(*) FROM (`+snapshotQuestions+`) q)
				AND NOT EXISTS (
					SELECT 1 FROM (`+snapshotQuestions+`) q
					LEFT JOIN LATERAL (
						SELECT count(*) AS copies, jsonb_agg(item)->0 AS item
						FROM jsonb_array_elements($2::jsonb) item
						WHERE item->>'questionId' = q.id::text
					) e ON true
					WHERE e.copies <> 1 OR CASE
						WHEN q.question->>'format' = 'code' THEN
							$3 <> '' OR e.item ? 'score' OR coalesce(jsonb_array_length(q.question->'test_cases'), 0) = 0
							OR (e.item->>'totalCases')::int IS DISTINCT FROM jsonb_array_length(q.question->'test_cases')
							OR NOT coalesce((e.item->>'passedCases')::int BETWEEN 0 AND (e.item->>'totalCases')::int, false)
						ELSE
							e.item ? 'passedCases' OR e.item ? 'totalCases'
							OR ((`+behavioralQuestion+`) = (e.item ? 'score'))
						END
				)
			RETURNING instance_id, status
		), mirrored AS (
			UPDATE assessment_instances i SET status = f.status FROM finalized f WHERE i.id = f.instance_id
			RETURNING i.id
		)
		SELECT EXISTS (SELECT 1 FROM mirrored) OR EXISTS (
			SELECT 1 FROM assessment_attempts a
			JOIN stage_executions se ON se.id = a.stage_execution_id
			JOIN pipeline_executions pe ON pe.id = se.pipeline_execution_id
			JOIN applications app ON app.id = pe.application_id
			JOIN job_postings j ON j.id = app.posting_id
			WHERE a.id = $1::uuid AND a.status = 'scored' AND a.final_evidence = $2::jsonb
				AND ($3 = '' OR j.company_id::text = $3)
		)
	`, attemptID, evidence, company).Scan(&finalized)
	if err != nil {
		return err
	}
	if !finalized {
		return ErrAssessmentConflict
	}
	return tx.Commit(ctx)
}

func (r *AssessmentRepository) RecordInterviewOutcome(ctx context.Context, interview, company, interviewer, outcome string) (json.RawMessage, error) {
	var result json.RawMessage
	err := r.pool.QueryRow(ctx, `
		UPDATE interviews i SET status = 'completed', outcome = $4,
			updated_at = CASE WHEN i.status = 'completed' THEN i.updated_at ELSE statement_timestamp() END
		FROM stage_executions se
		JOIN pipeline_executions pe ON pe.id = se.pipeline_execution_id
		JOIN applications a ON a.id = pe.application_id
		JOIN job_postings j ON j.id = a.posting_id
		WHERE i.id = $1::uuid AND i.stage_execution_id = se.id
			AND (($2 <> '' AND j.company_id::text = $2) OR EXISTS (
				SELECT 1 FROM interviewer_assignments ia
				WHERE ia.interview_id = i.id AND ia.interviewer_id::text = $3
			))
			AND ((i.status = 'confirmed' AND i.scheduled_end_at <= statement_timestamp() AND pe.status = 'active')
				OR (i.status = 'completed' AND i.outcome = $4))
		RETURNING jsonb_build_object('id', i.id, 'stageExecutionId', i.stage_execution_id,
			'scheduledStartAt', i.scheduled_start_at, 'scheduledEndAt', i.scheduled_end_at,
			'status', i.status, 'outcome', i.outcome, 'rescheduleCount', i.reschedule_count,
			'createdAt', i.created_at, 'updatedAt', i.updated_at)
	`, interview, company, interviewer, outcome).Scan(&result)
	return result, err
}

func (r *AssessmentRepository) ListInstances(ctx context.Context, candidate, application string) (json.RawMessage, error) {
	var result json.RawMessage
	err := r.pool.QueryRow(ctx, `
		SELECT coalesce(jsonb_agg(jsonb_build_object('id', i.id, 'applicationId', i.application_id,
			'stageId', i.stage_id, 'status', a.status, 'startedAt', a.started_at, 'deadline', a.deadline)
			ORDER BY a.deadline, i.id) FILTER (WHERE a.id IS NOT NULL), '[]'::jsonb)
		FROM applications app
		LEFT JOIN assessment_instances i ON i.application_id = app.id
		LEFT JOIN assessment_attempts a ON a.instance_id = i.id AND a.candidate_id = app.cand_id
		WHERE app.id = $2::uuid AND app.cand_id = $1::uuid
		GROUP BY app.id
	`, candidate, application).Scan(&result)
	return result, err
}

func (r *AssessmentRepository) ReviewEvaluation(ctx context.Context, attempt, company string) (json.RawMessage, error) {
	var result json.RawMessage
	err := r.pool.QueryRow(ctx, `
		SELECT jsonb_build_object('status', a.status, 'submittedAt', a.submitted_at,
			'finalScore', a.final_score, 'evaluatedAt', a.evaluated_at, 'evidence', a.final_evidence,
			'questions', coalesce((
				SELECT jsonb_agg(jsonb_build_object('id', q.id, 'moduleId', q.module_id,
					'prompt', q.question->'prompt', 'format', q.question->'format',
					'response', a.saved_responses->q.id::text) ORDER BY q.module_ordinal, q.question_ordinal)
				FROM (`+snapshotQuestions+`) q
			), '[]'::jsonb))
		FROM assessment_attempts a
		JOIN stage_executions se ON se.id = a.stage_execution_id
		JOIN pipeline_executions pe ON pe.id = se.pipeline_execution_id
		JOIN applications app ON app.id = pe.application_id
		JOIN job_postings j ON j.id = app.posting_id
		WHERE a.id = $1::uuid AND j.company_id::text = $2 AND a.submitted_at IS NOT NULL
	`, attempt, company).Scan(&result)
	return result, err
}

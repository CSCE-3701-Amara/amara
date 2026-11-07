package database

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/invw/domain"
	"github.com/jackc/pgx/v5"
)

type CodeRunner interface {
	Run(ctx context.Context, language, source, input string) (stdout string, succeeded bool, err error)
}

type codeGradingQuestion struct {
	ID       string
	Question domain.Question
	Response struct {
		Code     string
		Language string
	}
}

// SaveMixedEvidence records the complete non-code portion of a mixed attempt.
// Once accepted it is immutable, so a worker can safely combine it with runs
// made outside the database transaction.
func (r *AssessmentRepository) SaveMixedEvidence(ctx context.Context, attempt, company string, evidence []QuestionEvidence) error {
	if company == "" || len(evidence) == 0 {
		return ErrAssessmentConflict
	}
	for _, item := range evidence {
		if item.PassedCases != nil || item.TotalCases != nil ||
			(item.Score != nil && (*item.Score < 0 || *item.Score > 100)) {
			return ErrAssessmentConflict
		}
	}
	evidence = slices.Clone(evidence)
	slices.SortFunc(evidence, func(a, b QuestionEvidence) int {
		return strings.Compare(a.QuestionID, b.QuestionID)
	})
	tag, err := r.pool.Exec(ctx, `
		UPDATE assessment_attempts a SET non_code_evidence = $3::jsonb
		FROM stage_executions se
		JOIN pipeline_executions pe ON pe.id = se.pipeline_execution_id
		JOIN applications app ON app.id = pe.application_id
		JOIN job_postings j ON j.id = app.posting_id
		WHERE a.id = $1::uuid AND a.stage_execution_id = se.id AND j.company_id::text = $2
			AND a.submitted_at IS NOT NULL
			AND (pe.status = 'active' OR se.status IN ('passed', 'rejected'))
			AND (
				(a.status = 'pending_grading' AND a.non_code_evidence IS NULL)
				OR (a.status IN ('pending_grading', 'scored') AND a.non_code_evidence = $3::jsonb)
			)
			AND EXISTS (SELECT 1 FROM (`+snapshotQuestions+`) q WHERE q.question->>'format' = 'code')
			AND EXISTS (SELECT 1 FROM (`+snapshotQuestions+`) q
				WHERE q.question->>'format' <> 'code' AND NOT `+behavioralQuestion+`)
			AND NOT EXISTS (SELECT 1 FROM (`+snapshotQuestions+`) q
				WHERE q.question->>'format' = 'code'
					AND coalesce(jsonb_array_length(q.question->'test_cases'), 0) = 0)
			AND jsonb_array_length($3::jsonb) = (SELECT count(*) FROM (`+snapshotQuestions+`) q
				WHERE q.question->>'format' <> 'code')
			AND NOT EXISTS (
				SELECT 1 FROM (`+snapshotQuestions+`) q
				WHERE q.question->>'format' <> 'code' AND 1 <> (
					SELECT count(*) FROM jsonb_array_elements($3::jsonb) item
					WHERE item->>'questionId' = q.id::text
						AND ((item ? 'score') <> (`+behavioralQuestion+`))
				)
			)
	`, attempt, company, evidence)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrAssessmentConflict
	}
	return nil
}

// GradeNextCodeAttempt claims one submitted attempt. The retry timestamp also
// acts as a lease so concurrent workers do not grade the same attempt.
func (r *AssessmentRepository) GradeNextCodeAttempt(ctx context.Context, runner CodeRunner) (bool, error) {
	var attempt string
	var lease time.Time
	var questions []codeGradingQuestion
	var evidence []QuestionEvidence
	err := r.pool.QueryRow(ctx, `
		WITH next AS (
			SELECT a.id
			FROM assessment_attempts a
			JOIN stage_executions se ON se.id = a.stage_execution_id
			JOIN pipeline_executions pe ON pe.id = se.pipeline_execution_id
			WHERE a.status = 'pending_grading' AND a.submitted_at IS NOT NULL
				AND (a.code_grading_retry_at IS NULL OR a.code_grading_retry_at <= statement_timestamp())
				AND (pe.status = 'active' OR se.status IN ('passed', 'rejected'))
				AND EXISTS (SELECT 1 FROM (`+snapshotQuestions+`) q WHERE q.question->>'format' = 'code')
				AND NOT EXISTS (
					SELECT 1 FROM (`+snapshotQuestions+`) q
					WHERE (q.question->>'format' = 'code'
						AND coalesce(jsonb_array_length(q.question->'test_cases'), 0) = 0)
						OR (q.question->>'format' <> 'code' AND NOT `+behavioralQuestion+`
							AND a.non_code_evidence IS NULL)
				)
			ORDER BY a.code_grading_retry_at NULLS FIRST, a.submitted_at, a.id
			LIMIT 1 FOR UPDATE OF a SKIP LOCKED
		)
		UPDATE assessment_attempts a
		SET code_grading_retry_at = statement_timestamp() + interval '5 minutes'
		FROM next WHERE a.id = next.id
		RETURNING a.id::text, a.code_grading_retry_at, a.non_code_evidence, (
			SELECT jsonb_agg(jsonb_build_object('id', q.id, 'question', q.question,
				'response', CASE WHEN q.question->>'format' = 'code' THEN a.saved_responses->q.id::text END)
				ORDER BY q.module_ordinal, q.question_ordinal)
			FROM (`+snapshotQuestions+`) q
		)`).Scan(&attempt, &lease, &evidence, &questions)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	hasNonCodeEvidence := len(evidence) > 0
	for _, q := range questions {
		if q.Question.Format != "code" && hasNonCodeEvidence {
			continue
		}
		item := QuestionEvidence{QuestionID: q.ID}
		if q.Question.Format == "code" {
			passed, total := 0, len(q.Question.TestCases)
			if strings.TrimSpace(q.Response.Code) != "" {
				for _, test := range q.Question.TestCases {
					// Renew before starting a case if the lease may expire during it.
					if time.Until(lease) < time.Minute {
						err := r.pool.QueryRow(ctx, `
							UPDATE assessment_attempts
							SET code_grading_retry_at = statement_timestamp() + interval '5 minutes'
							WHERE id = $1::uuid AND status = 'pending_grading'
								AND code_grading_retry_at = $2 AND code_grading_retry_at > statement_timestamp()
							RETURNING code_grading_retry_at
						`, attempt, lease).Scan(&lease)
						if err != nil {
							return true, fmt.Errorf("renew code grading lease for attempt %s: %w", attempt, err)
						}
					}
					output, succeeded, err := runner.Run(ctx, q.Response.Language, q.Response.Code, test.Input)
					if err != nil {
						return true, fmt.Errorf("grade attempt %s, question %s: %w", attempt, q.ID, err)
					}
					if succeeded && normalizeCodeOutput(output) == normalizeCodeOutput(test.ExpectedOutput) {
						passed++
					}
				}
			}
			item.PassedCases, item.TotalCases = &passed, &total
		}
		evidence = append(evidence, item)
	}
	if err := r.FinalizeEvidence(ctx, attempt, "", evidence); err != nil {
		return true, fmt.Errorf("finalize code grading for attempt %s: %w", attempt, err)
	}
	return true, nil
}

func normalizeCodeOutput(value string) string {
	return strings.TrimRight(strings.ReplaceAll(value, "\r\n", "\n"), "\n")
}

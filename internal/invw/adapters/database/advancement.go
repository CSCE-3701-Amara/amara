package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CSCE-3701-Amara/amara/internal/invw/domain"
)

var ErrDecisionConflict = errors.New("stage is not eligible for this manual decision")

type advancementStage struct {
	application, pipeline, version, stage, candidate                                    string
	applicationStatus, pipelineStatus, status, rule, below, attempt, outcome, interview string
	threshold, score                                                                    *int
	submitted, behavioral, evaluated                                                    bool
	content                                                                             []struct {
		ModuleID  string `json:"module_id"`
		Questions []domain.Question
	}
	configuration struct {
		Modules []struct {
			ID            string
			Configuration struct{ Definition domain.Definition }
		}
	}
	evidence    []QuestionEvidence
	questionIDs map[string]string
}

func readAdvancement(ctx context.Context, tx pgx.Tx, id, company, stage string) (advancementStage, error) {
	var s advancementStage
	// The caller already locked the application. UC-010 commits evidence and
	// the outcome together, and a submitted attempt cannot become missed.
	err := tx.QueryRow(ctx, `
 SELECT a.id::text, pe.id::text, pe.pipeline_version_id::text, ps.id::text, a.cand_id::text,
 coalesce(a.status,''), pe.status, se.status, ps.advancement_rule, ps.below_threshold_action, ps.score_threshold,
 coalesce(t.id::text,''), coalesce(t.status,''), coalesce(i.outcome,''),
 t.submitted_at IS NOT NULL, coalesce((ps.interview_configuration->>'behavioral')::boolean,false),
 t.evaluated_at IS NOT NULL, coalesce(t.final_score,i.score),
 coalesce(t.selected_content,'[]'), coalesce(t.configuration_snapshot,'{}'), coalesce(t.final_evidence,'[]'),
 coalesce((SELECT jsonb_object_agg(q.module_ordinal || ':' || q.question_ordinal, q.id)
   FROM (SELECT t.id, t.selected_content) a
   CROSS JOIN LATERAL (`+snapshotQuestions+`) q), '{}')
 FROM stage_executions se
 JOIN pipeline_executions pe ON pe.id = se.pipeline_execution_id
 JOIN applications a ON a.id = pe.application_id
 JOIN job_postings j ON j.id = a.posting_id
 JOIN pipeline_stages ps ON ps.id = se.stage_id AND ps.version_id = pe.pipeline_version_id
 LEFT JOIN assessment_attempts t ON t.stage_execution_id = se.id
 -- UC-010 consumes interview records once scheduling and scorecard submission
 -- are connected. Until then an available interview stage awaits those records.
 LEFT JOIN interviews i ON i.stage_execution_id = se.id AND i.status = 'completed'
 WHERE se.id = $1::text::uuid AND ($2 = '' OR j.company_id::text = $2)
 AND ($3 = '' OR se.stage_id::text = $3)`, id, company, stage).Scan(&s.application, &s.pipeline, &s.version, &s.stage, &s.candidate,
		&s.applicationStatus, &s.pipelineStatus, &s.status, &s.rule, &s.below, &s.threshold,
		&s.attempt, &s.outcome, &s.interview, &s.submitted, &s.behavioral, &s.evaluated, &s.score,
		&s.content, &s.configuration, &s.evidence, &s.questionIDs)
	return s, err
}

type stageDecision struct {
	status, kind, reason string
	score                *int
}

func (s advancementStage) automaticDecision() stageDecision {
	if s.outcome == "missed" || s.interview == "candidate_absent" {
		return stageDecision{status: "rejected", kind: "missed", reason: "Missed module"}
	}
	if !s.submitted && s.interview != "held" {
		return stageDecision{}
	}
	for _, module := range s.configuration.Modules {
		definition := module.Configuration.Definition
		if definition.Type == domain.Behavioral || definition.Behavioral {
			s.behavioral = true
		}
	}
	if s.behavioral {
		return stageDecision{status: "held"}
	}

	// The external evaluator must call UC-010 FinalizeEvidence with trusted
	// per-question results. Until it is connected, coding attempts stay pending.
	coding, valid, passed := s.codingThreshold()
	if coding && !valid {
		return stageDecision{status: "pending_grading"}
	}
	if s.rule == "manual" {
		return stageDecision{status: "held", score: s.score}
	}
	if coding {
		return s.thresholdDecision(passed)
	}
	if s.rule == "threshold" {
		// Live scores are populated by UC-010 from submitted human scorecards.
		// The scorecard submission route is still unwired.
		if s.score == nil || s.threshold == nil || (s.attempt != "" && !s.evaluated) {
			return stageDecision{status: "pending_grading"}
		}
		return s.thresholdDecision(*s.score >= *s.threshold)
	}
	return stageDecision{status: "passed", kind: "completion", score: s.score}
}

func (s advancementStage) thresholdDecision(passed bool) stageDecision {
	if passed {
		return stageDecision{status: "passed", kind: "threshold", score: s.score}
	}
	if s.below == "reject" {
		return stageDecision{status: "rejected", kind: "threshold", reason: "Below threshold", score: s.score}
	}
	return stageDecision{status: "held", score: s.score}
}

func (s advancementStage) codingThreshold() (coding, valid, passed bool) {
	evidence := make(map[string]QuestionEvidence, len(s.evidence))
	for _, result := range s.evidence {
		evidence[result.QuestionID] = result
	}
	passed = true
	for index, module := range s.configuration.Modules {
		definition := module.Configuration.Definition
		if definition.Type != domain.Coding {
			continue
		}
		coding = true
		if !s.evaluated || s.outcome != "scored" || definition.MinimumTestsPassed == nil ||
			index >= len(s.content) || s.content[index].ModuleID != module.ID {
			return true, false, false
		}
		count := 0
		for position, question := range s.content[index].Questions {
			result := evidence[s.questionIDs[fmt.Sprintf("%d:%d", index+1, position+1)]]
			if question.MinimumCasesPassed == nil || result.PassedCases == nil || result.TotalCases == nil ||
				*result.TotalCases != len(question.TestCases) || *result.PassedCases < 0 || *result.PassedCases > *result.TotalCases {
				return true, false, false
			}
			if *result.PassedCases >= *question.MinimumCasesPassed {
				count++
			}
		}
		passed = passed && count >= *definition.MinimumTestsPassed
	}
	return coding, true, passed
}

func applyAdvancement(ctx context.Context, tx pgx.Tx, id, actor string, s advancementStage, d stageDecision) error {
	if d.status == "" {
		return nil
	}
	if d.status == "held" || d.status == "pending_grading" {
		_, err := tx.Exec(ctx, `UPDATE stage_executions SET status=$3, score=$2
			WHERE id=$1::text::uuid AND (status,score) IS DISTINCT FROM ($3::text,$2::integer)`, id, d.score, d.status)
		return err
	}
	_, err := tx.Exec(ctx, `
 UPDATE stage_executions SET status=$2, score=$3, score_at_decision=$3,
 completed_at=now(), decision_type=$4, decided_by=nullif($5::text,'')::uuid, decided_at=now(), decision_reason=nullif($6::text,''),
 decision_history=decision_history || jsonb_build_array(jsonb_build_object(
  'candidateId',$7::text,'stageExecutionId',$1::text,'decision',CASE WHEN $2='passed' THEN 'pass' ELSE 'reject' END,
  'decisionType',$4::text,'decidedBy',nullif($5::text,''),'decidedAt',now(),'reason',nullif($6::text,''),'scoreAtDecision',$3::integer))
 WHERE id=$1::text::uuid`, id, d.status, d.score, d.kind, actor, d.reason, s.candidate)
	if err != nil {
		return err
	}
	if d.status == "rejected" {
		// Persist the feedback work item on the existing pipeline record. Feedback
		// generation, approval and delivery are pending UC-INVW-009 integration.
		_, err = tx.Exec(ctx, `
   UPDATE pipeline_executions SET status='rejected',completed_at=now(),updated_at=now(),
    rejection_feedback=coalesce(rejection_feedback,jsonb_build_object('status','pending','stageExecutionId',$2::text))
   WHERE id=$1::text::uuid`, s.pipeline, id)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `UPDATE stage_executions SET status='cancelled',completed_at=now()
   WHERE pipeline_execution_id=$1::text::uuid AND status IN ('locked','available','in_progress','pending_grading','held')`, s.pipeline)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `UPDATE applications SET status='rejected',final_decision=false,updated_at=now() WHERE id=$1::text::uuid`, s.application)
		return err
	}

	_, err = tx.Exec(ctx, `
 INSERT INTO candidate_amarat (cand_id,amara_title,amara_description,module_id,company_id,application_id,
  source_stage_execution_id,definition_snapshot,score,pass_decision)
 SELECT a.cand_id, amara->>'name', amara->>'description',tc.test_module_id,j.company_id,a.id,se.id,
  amara,se.score,se.decision_history->-1
 FROM stage_executions se
 JOIN pipeline_executions pe ON pe.id=se.pipeline_execution_id
 JOIN applications a ON a.id=pe.application_id
 JOIN job_postings j ON j.id=a.posting_id
 JOIN test_configurations tc ON tc.stage_id=se.stage_id
 CROSS JOIN LATERAL jsonb_array_elements(coalesce(tc.configuration#>'{definition,amarat}','[]')) amara
 WHERE se.id=$1::text::uuid`, id)
	if err != nil {
		return err
	}

	// Older DAG versions use explicit links. Current linear versions store order.
	_, err = tx.Exec(ctx, `
 WITH links AS (
  SELECT from_stage_id AS source,to_stage_id AS target FROM stage_transitions WHERE pipeline_version_id=$2::text::uuid
  UNION ALL
  SELECT p.id,n.id FROM pipeline_stages p JOIN pipeline_stages n
   ON n.version_id=p.version_id AND n.stage_order=p.stage_order+1
  WHERE p.version_id=$2::text::uuid AND NOT EXISTS (
   SELECT 1 FROM stage_transitions WHERE pipeline_version_id=p.version_id)
 ), ready AS (
  SELECT n.id,n.scheduled_start_at,n.scheduled_end_at FROM pipeline_stages n
  WHERE n.version_id=$2::text::uuid
   AND EXISTS (SELECT 1 FROM links WHERE source=$3::text::uuid AND target=n.id)
   AND NOT EXISTS (
    SELECT 1 FROM links l WHERE l.target=n.id AND NOT EXISTS (
     SELECT 1 FROM stage_executions se WHERE se.pipeline_execution_id=$1::text::uuid
      AND se.stage_id=l.source AND se.status='passed'))
 )
 INSERT INTO stage_executions (id,pipeline_execution_id,stage_id,status,available_at,deadline)
 SELECT gen_random_uuid(),$1::text::uuid,id,'available',greatest(now(),scheduled_start_at),scheduled_end_at FROM ready
 ON CONFLICT (pipeline_execution_id,stage_id) DO UPDATE SET status='available',
  available_at=excluded.available_at,deadline=excluded.deadline WHERE stage_executions.status='locked'`,
		s.pipeline, s.version, s.stage)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
 WITH finished AS (
  UPDATE pipeline_executions pe SET status='completed',completed_at=now(),updated_at=now()
  WHERE pe.id=$1::text::uuid AND NOT EXISTS (
   SELECT 1 FROM pipeline_stages ps WHERE ps.version_id=pe.pipeline_version_id AND NOT EXISTS (
    SELECT 1 FROM stage_executions se WHERE se.pipeline_execution_id=pe.id AND se.stage_id=ps.id AND se.status='passed'))
  RETURNING application_id
 )
 UPDATE applications SET status='finalist',updated_at=now() WHERE id IN (SELECT application_id FROM finished)`, s.pipeline)
	return err
}

func (s advancementStage) active() bool {
	if s.applicationStatus != "admitted" || s.pipelineStatus != "active" {
		return false
	}
	switch s.status {
	case "available", "in_progress", "pending_grading", "held", "missed":
		return true
	default:
		return false
	}
}

// ReconcileAdvancement consumes UC-010 outcomes through the worker's normal sweep.
func ReconcileAdvancement(ctx context.Context, pool *pgxpool.Pool) error {
	rows, err := pool.Query(ctx, `
 SELECT se.id::text FROM stage_executions se
 JOIN pipeline_executions pe ON pe.id=se.pipeline_execution_id AND pe.status='active'
 LEFT JOIN assessment_attempts t ON t.stage_execution_id=se.id
 LEFT JOIN interviews i ON i.stage_execution_id=se.id AND i.status='completed'
 WHERE se.status IN ('available','in_progress','pending_grading','held','missed')
 AND (t.status='missed' OR t.submitted_at IS NOT NULL OR i.outcome IN ('held','candidate_absent'))
 ORDER BY (coalesce(t.status='missed',false) OR coalesce(i.outcome='candidate_absent',false)) DESC,se.id`)
	if err != nil {
		return err
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}
	var result error
	for _, id := range ids {
		_, err := decideStages(ctx, pool, []string{id}, "", "", "", "", "")
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			result = errors.Join(result, err)
		}
	}
	return result
}

func (r *AssessmentRepository) DecideStages(ctx context.Context, ids []string, company, stage, decision, reason, expected string) ([]json.RawMessage, error) {
	if company == "" || (decision != "pass" && decision != "reject") {
		return nil, ErrDecisionConflict
	}
	return decideStages(ctx, r.pool, ids, company, stage, decision, reason, expected)
}

func decideStages(ctx context.Context, pool *pgxpool.Pool, ids []string, company, stage, decision, reason, expected string) ([]json.RawMessage, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Serialize decisions on the same application, including parallel stages.
	// Acquire locks in application order for bulk requests. NO KEY UPDATE is
	// compatible with foreign-key checks when the worker creates attempts.
	if _, err = tx.Exec(ctx, `
 SELECT a.id FROM applications a JOIN pipeline_executions pe ON pe.application_id=a.id
 JOIN stage_executions se ON se.pipeline_execution_id=pe.id
 JOIN job_postings j ON j.id=a.posting_id
 WHERE se.id=ANY($1::text[]::uuid[]) AND ($2='' OR j.company_id::text=$2)
 ORDER BY a.id FOR NO KEY UPDATE OF a`, ids, company); err != nil {
		return nil, err
	}

	results := make([]json.RawMessage, 0, len(ids))
	for _, id := range ids {
		s, err := readAdvancement(ctx, tx, id, company, stage)
		if err != nil {
			return nil, err
		}
		if s.active() {
			d := s.automaticDecision()
			actor := ""
			if decision != "" && d.status != "passed" && d.status != "rejected" {
				if d.status != "held" || (expected != "" && expected != s.status && expected != d.status) {
					return nil, ErrDecisionConflict
				}
				d.kind, d.reason, actor = "manual", reason, company
				d.status = "rejected"
				if decision == "pass" {
					d.status = "passed"
				}
			}
			if err = applyAdvancement(ctx, tx, id, actor, s, d); err != nil {
				return nil, err
			}
		} else if decision != "" && s.status != "passed" && s.status != "rejected" {
			return nil, ErrDecisionConflict
		}

		if decision != "" {
			var result json.RawMessage
			err = tx.QueryRow(ctx, `
    SELECT jsonb_build_object('id',id,'pipelineExecutionId',pipeline_execution_id,'stageId',stage_id,
     'status',status,'availableAt',available_at,'deadline',deadline,'startedAt',started_at,
     'completedAt',completed_at,'score',score,'decisionHistory',decision_history)
    FROM stage_executions WHERE id=$1::text::uuid`, id).Scan(&result)
			if err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	}
	return results, tx.Commit(ctx)
}

package database

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrThresholdConflict = errors.New("threshold or affected candidates changed")

func (r *AssessmentRepository) PipelineDashboard(ctx context.Context, posting, company string) (json.RawMessage, error) {
	var dashboard json.RawMessage
	err := r.pool.QueryRow(ctx, `
 WITH published AS (
  SELECT v.id, j.finalist_target FROM pipeline_versions v
  JOIN pipelines p ON p.id=v.pipeline_id
  JOIN job_postings j ON j.id=p.posting_id
  WHERE j.id=$1::uuid AND j.company_id=$2::uuid AND j.deleted_at IS NULL
  AND (v.id=p.published_version_id OR (p.published_version_id IS NULL AND v.is_published))
 ), candidates AS (
  SELECT pe.id, pe.anonymous_identifier AS candidate, pe.status, a.nice_to_have_score
  FROM published v JOIN pipeline_executions pe ON pe.pipeline_version_id=v.id
  JOIN applications a ON a.id=pe.application_id
 ), scores AS (
  SELECT se.id, se.stage_id, se.pipeline_execution_id, se.status,
   coalesce(t.final_score,se.score,i.score) AS score,
   c.candidate, c.nice_to_have_score,
   c.status='active' AND se.status IN ('available','in_progress','pending_grading','held') AS active
  FROM candidates c JOIN stage_executions se ON se.pipeline_execution_id=c.id
  LEFT JOIN assessment_attempts t ON t.stage_execution_id=se.id
  LEFT JOIN interviews i ON i.stage_execution_id=se.id AND i.status='completed'
 ), percentiles AS (
  SELECT id, round(100*cume_dist() OVER (PARTITION BY stage_id ORDER BY score)) AS percentile
  FROM scores WHERE score IS NOT NULL
 ), stages AS (
  SELECT ps.id, ps.name, ps.stage_order, ps.advancement_rule AS "rule", ps.score_threshold AS threshold,
   count(s.id) FILTER (WHERE s.active) AS active,
   count(s.id) FILTER (WHERE s.status='passed') AS passed,
   count(s.id) FILTER (WHERE s.status IN ('rejected','missed')) AS rejected,
   coalesce(jsonb_agg(jsonb_build_object('id',s.id,'candidate',s.candidate,'status',s.status,
     'active',s.active,'score',s.score,'niceToHaveScore',s.nice_to_have_score,'percentile',p.percentile)
     ORDER BY s.score DESC NULLS LAST,s.candidate) FILTER (WHERE s.id IS NOT NULL),'[]') AS candidates
  FROM published v JOIN pipeline_stages ps ON ps.version_id=v.id
  LEFT JOIN scores s ON s.stage_id=ps.id AND s.status <> 'locked'
  LEFT JOIN percentiles p ON p.id=s.id
  GROUP BY ps.id
 )
 SELECT jsonb_build_object('finalistTarget',v.finalist_target,
  'active',(SELECT count(*) FROM candidates WHERE status='active'),
  'links',(SELECT coalesce(jsonb_agg(jsonb_build_object('from',from_stage_id,'to',to_stage_id)),'[]')
   FROM stage_transitions WHERE pipeline_version_id=v.id),
  'stages',(SELECT coalesce(jsonb_agg(to_jsonb(stages) ORDER BY stage_order,id),'[]') FROM stages))
 FROM published v`, posting, company).Scan(&dashboard)
	return dashboard, err
}

type ThresholdOutcome struct {
	Candidate string `json:"candidate"`
	Status    string `json:"status"`
	Score     int    `json:"score"`
}

func (r *AssessmentRepository) ChangeStageThreshold(ctx context.Context, stage, company string, threshold, expected int, apply bool, preview []ThresholdOutcome) ([]ThresholdOutcome, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var previous int
	err = tx.QueryRow(ctx, `
 SELECT ps.score_threshold FROM pipeline_stages ps
 JOIN pipeline_versions v ON v.id=ps.version_id
 JOIN pipelines p ON p.id=v.pipeline_id
 JOIN job_postings j ON j.id=p.posting_id
 WHERE ps.id=$1::uuid AND j.company_id=$2::uuid AND j.deleted_at IS NULL
 AND (v.id=p.published_version_id OR (p.published_version_id IS NULL AND v.is_published))
 AND ps.advancement_rule='threshold'
 FOR NO KEY UPDATE OF ps`, stage, company).Scan(&previous)
	if err != nil {
		return nil, err
	}
	if previous != expected {
		return nil, ErrThresholdConflict
	}
	rows, err := tx.Query(ctx, `
 SELECT se.id::text,pe.anonymous_identifier::text
 FROM stage_executions se JOIN pipeline_executions pe ON pe.id=se.pipeline_execution_id
 JOIN applications a ON a.id=pe.application_id
 WHERE se.stage_id=$1::uuid AND pe.status='active' AND a.status='admitted'
 AND se.status IN ('available','in_progress','pending_grading','held')
 ORDER BY a.id FOR NO KEY UPDATE OF a`, stage)
	if err != nil {
		return nil, err
	}
	type waiting struct{ ID, Candidate string }
	items, err := pgx.CollectRows(rows, pgx.RowToStructByPos[waiting])
	if err != nil {
		return nil, err
	}
	outcomes := make([]ThresholdOutcome, 0)
	for _, item := range items {
		s, err := readAdvancement(ctx, tx, item.ID, company, stage)
		if err != nil {
			return nil, err
		}
		if !s.active() || s.score == nil {
			continue
		}
		before := s.automaticDecision()
		s.threshold = &threshold
		after := s.automaticDecision()
		if after.status == before.status || after.status == s.status || after.score == nil {
			continue
		}
		outcome := ThresholdOutcome{Candidate: item.Candidate, Status: after.status, Score: *after.score}
		if apply {
			if len(outcomes) >= len(preview) || preview[len(outcomes)] != outcome {
				return nil, ErrThresholdConflict
			}
			if err = applyAdvancement(ctx, tx, item.ID, company, s, after); err != nil {
				return nil, err
			}
		}
		outcomes = append(outcomes, outcome)
	}
	if !apply {
		return outcomes, nil
	}
	if len(preview) != len(outcomes) {
		return nil, ErrThresholdConflict
	}
	if threshold != previous {
		if _, err = tx.Exec(ctx, `UPDATE pipeline_stages SET score_threshold=$2 WHERE id=$1::uuid`, stage, threshold); err != nil {
			return nil, err
		}
		if _, err = tx.Exec(ctx, `
 INSERT INTO pipeline_threshold_changes(stage_id,changed_by,previous_threshold,new_threshold)
 VALUES ($1::uuid,$2::uuid,$3,$4)`, stage, company, previous, threshold); err != nil {
			return nil, err
		}
	}
	return outcomes, tx.Commit(ctx)
}

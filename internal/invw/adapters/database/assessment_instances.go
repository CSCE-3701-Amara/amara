package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SpawnAssessmentInstances admits screened-in applications to their published
// pipeline and snapshots the content of each available assessment.
func SpawnAssessmentInstances(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		WITH entry_stages AS (
			SELECT s.* FROM pipeline_stages s
			WHERE CASE WHEN EXISTS (
				SELECT 1 FROM stage_transitions WHERE pipeline_version_id = s.version_id
			) THEN NOT EXISTS (
				SELECT 1 FROM stage_transitions WHERE pipeline_version_id = s.version_id AND to_stage_id = s.id
			) ELSE s.stage_order = 1 END
		), admitted AS (
			INSERT INTO pipeline_executions (id, application_id, pipeline_version_id)
			SELECT gen_random_uuid(), a.id, v.id
			FROM applications a
			JOIN pipelines p ON p.posting_id = a.posting_id
			JOIN pipeline_versions v ON v.id = p.published_version_id AND v.pipeline_id = p.id
			WHERE a.status = 'admitted' AND a.screened_out IS FALSE
				AND EXISTS (SELECT 1 FROM entry_stages s WHERE s.version_id = v.id)
			ON CONFLICT (application_id) DO NOTHING
			RETURNING id, pipeline_version_id, started_at
		)
		INSERT INTO stage_executions (id, pipeline_execution_id, stage_id, status, available_at, deadline)
		SELECT gen_random_uuid(), a.id, s.id, 'available',
			greatest(a.started_at, s.scheduled_start_at), s.scheduled_end_at
		FROM admitted a JOIN entry_stages s ON s.version_id = a.pipeline_version_id
	`); err != nil {
		return fmt.Errorf("admit applications to published pipelines: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		WITH eligible AS MATERIALIZED (
			SELECT se.id AS stage_execution_id, pe.application_id, a.cand_id, se.stage_id,
				least(se.deadline, s.scheduled_end_at, se.available_at + (
					SELECT min(tc.completion_window) FROM test_configurations tc WHERE tc.stage_id = s.id
				)) AS deadline
			FROM stage_executions se
			JOIN pipeline_executions pe ON pe.id = se.pipeline_execution_id AND pe.status = 'active'
			JOIN pipeline_stages s ON s.id = se.stage_id AND s.version_id = pe.pipeline_version_id
			JOIN applications a ON a.id = pe.application_id AND a.status = 'admitted'
			WHERE se.status = 'available' AND s.module_type = 'assessment'
				AND s.scheduled_start_at IS NOT NULL AND s.scheduled_end_at IS NOT NULL
				AND se.available_at IS NOT NULL
				AND greatest(se.available_at, s.scheduled_start_at) <= CURRENT_TIMESTAMP
				AND EXISTS (SELECT 1 FROM test_configurations tc WHERE tc.stage_id = s.id)
				AND NOT EXISTS (SELECT 1 FROM assessment_attempts t WHERE t.stage_execution_id = se.id)
			FOR UPDATE OF se SKIP LOCKED
		), created AS (
			INSERT INTO assessment_instances (id, application_id, stage_id, status, deadline)
			SELECT gen_random_uuid(), application_id, stage_id,
				CASE WHEN deadline <= CURRENT_TIMESTAMP THEN 'missed' ELSE 'available' END,
				deadline AT TIME ZONE 'UTC'
			FROM eligible
			ON CONFLICT (application_id, stage_id) DO UPDATE SET
				status = EXCLUDED.status, deadline = EXCLUDED.deadline
			WHERE assessment_instances.started_at IS NULL
				AND assessment_instances.status IN ('available', 'missed')
			RETURNING id, application_id, stage_id, status
		)
		INSERT INTO assessment_attempts
			(id, stage_execution_id, candidate_id, instance_id, deadline, status, selected_content, configuration_snapshot)
		SELECT gen_random_uuid(), e.stage_execution_id, e.cand_id, i.id, e.deadline,
			i.status, content.questions, jsonb_build_object('modules', content.modules)
		FROM created i
		JOIN eligible e ON e.application_id = i.application_id AND e.stage_id = i.stage_id
		CROSS JOIN LATERAL (
			SELECT jsonb_agg(module_content.questions ORDER BY (tc.configuration->>'position')::int, tc.id) AS questions,
				jsonb_agg(jsonb_build_object('id', tc.test_module_id, 'active_time_limit', tc.active_time_limit,
					'configuration', tc.configuration) ORDER BY (tc.configuration->>'position')::int, tc.id) AS modules
			FROM test_configurations tc
			CROSS JOIN LATERAL (
				SELECT jsonb_build_object('module_id', tc.test_module_id, 'questions',
					coalesce(tc.configuration #> '{definition,questions}', '[]'::jsonb) ||
					coalesce((
						SELECT jsonb_agg(p.question ORDER BY p.position)
						FROM (
							SELECT m.authoring_definition->'questions'->((ref.item->>'position')::int) AS question,
								ref.position, row_number() OVER (
									PARTITION BY m.authoring_definition->'questions'->((ref.item->>'position')::int)->>'difficulty'
									ORDER BY random()
								) AS rank
							FROM jsonb_array_elements(coalesce(tc.configuration #> '{definition,library,problems}',
								tc.configuration #> '{definition,library,pool}', '[]'::jsonb))
								WITH ORDINALITY AS ref(item, position)
							JOIN modules m ON m.id = (ref.item->>'module_id')::uuid
						) p
						WHERE tc.configuration #> '{definition,library,difficulty_mix}' IS NULL
							OR p.rank <= (tc.configuration #>> ARRAY['definition', 'library', 'difficulty_mix', p.question->>'difficulty'])::int
					), '[]'::jsonb)
				) AS questions
			) module_content
			WHERE tc.stage_id = i.stage_id
		) content
	`); err != nil {
		return fmt.Errorf("spawn assessment instances: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		WITH expired AS (
			UPDATE assessment_attempts a
			SET status = CASE WHEN started_at IS NULL THEN 'missed'
				WHEN `+behavioralOnly+` THEN 'submitted' ELSE 'pending_grading' END,
				submitted_at = CASE WHEN started_at IS NOT NULL THEN deadline END,
				auto_submitted = started_at IS NOT NULL
			WHERE a.status IN ('available', 'in_progress') AND a.deadline <= statement_timestamp()
				AND EXISTS (
					SELECT 1 FROM stage_executions se
					JOIN pipeline_executions pe ON pe.id = se.pipeline_execution_id
					WHERE se.id = a.stage_execution_id AND pe.status = 'active'
						AND se.status IN ('available', 'in_progress')
				)
			RETURNING instance_id, status
		)
		UPDATE assessment_instances i SET status = e.status FROM expired e WHERE i.id = e.instance_id
	`); err != nil {
		return fmt.Errorf("record expired assessment outcomes: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		UPDATE interviews i SET score = (
			SELECT round((avg(r.rating::numeric) - 1) * 25)::int
			FROM interview_feedback f
			JOIN interviewer_assignments ia USING (interview_id, interviewer_id)
			CROSS JOIN LATERAL jsonb_each(f.ratings) r(criterion, rating)
			WHERE f.interview_id = i.id AND f.submitted_at IS NOT NULL AND ia.scorecard_status = 'submitted'
		), updated_at = statement_timestamp()
		FROM stage_executions se JOIN pipeline_stages s ON s.id = se.stage_id
		WHERE i.stage_execution_id = se.id AND i.status = 'completed' AND i.outcome = 'held'
			AND i.score IS NULL AND (s.interview_configuration->'behavioral' = 'true'::jsonb) IS NOT TRUE
			AND jsonb_array_length(i.rubric_snapshot) > 0
			AND EXISTS (SELECT 1 FROM interviewer_assignments ia
				WHERE ia.interview_id = i.id AND ia.scorecard_status = 'submitted')
			AND NOT EXISTS (
				SELECT 1 FROM interviewer_assignments ia
				LEFT JOIN interview_feedback f USING (interview_id, interviewer_id)
				WHERE ia.interview_id = i.id AND ia.scorecard_status <> 'waived'
					AND (ia.scorecard_status <> 'submitted' OR f.submitted_at IS NULL
						OR nullif(btrim(f.recommendation), '') IS NULL
						OR (SELECT count(*) FROM jsonb_each(f.ratings)) <> jsonb_array_length(i.rubric_snapshot)
						OR EXISTS (SELECT 1 FROM jsonb_each(f.ratings) r(criterion, rating)
							WHERE CASE WHEN jsonb_typeof(r.rating) = 'number' THEN
								r.rating::numeric NOT IN (1, 2, 3, 4, 5) ELSE true END))
			)
	`); err != nil {
		return fmt.Errorf("record completed interview scores: %w", err)
	}
	return tx.Commit(ctx)
}

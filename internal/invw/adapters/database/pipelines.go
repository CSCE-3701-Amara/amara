package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CSCE-3701-Amara/amara/internal/invw/domain"
)

type PipelineRepository struct {
	pool *pgxpool.Pool
}

func NewPipelineRepository(pool *pgxpool.Pool) *PipelineRepository {
	return &PipelineRepository{pool: pool}
}

func (r *PipelineRepository) Get(ctx context.Context, company, posting string) (domain.PipelineDraft, error) {
	var draft domain.PipelineDraft
	err := r.pool.QueryRow(ctx, `
		SELECT p.id::text, p.posting_id::text, v.id::text, v.version_number
		FROM pipelines p JOIN job_postings j ON j.id = p.posting_id
		JOIN LATERAL (
			SELECT id, version_number FROM pipeline_versions
			WHERE pipeline_id = p.id ORDER BY version_number DESC LIMIT 1
		) v ON true
		WHERE p.posting_id = $1::text::uuid AND j.company_id = $2::text::uuid`,
		posting, company).Scan(&draft.ID, &draft.PostingID, &draft.VersionID, &draft.VersionNumber)
	return draft, pipelineError(err)
}

func (r *PipelineRepository) Create(ctx context.Context, company, posting string) (domain.PipelineDraft, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.PipelineDraft{}, err
	}
	defer tx.Rollback(ctx)

	var attached bool
	err = tx.QueryRow(ctx, `
		SELECT pipeline_id IS NOT NULL FROM job_postings
		WHERE id = $1::text::uuid AND company_id = $2::text::uuid FOR UPDATE`, posting, company).Scan(&attached)
	if err != nil {
		return domain.PipelineDraft{}, pipelineError(err)
	}
	if attached {
		return domain.PipelineDraft{}, domain.ErrPipelineExists
	}

	draft := domain.PipelineDraft{PostingID: posting, VersionNumber: 1}
	err = tx.QueryRow(ctx, `
		INSERT INTO pipelines (id, posting_id, status) VALUES (gen_random_uuid(), $1::text::uuid, 'draft')
		ON CONFLICT (posting_id) DO NOTHING RETURNING id::text`, posting).Scan(&draft.ID)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.PipelineDraft{}, domain.ErrPipelineExists
	}
	if err != nil {
		return domain.PipelineDraft{}, err
	}
	err = tx.QueryRow(ctx, `
		INSERT INTO pipeline_versions (id, pipeline_id, version_number, is_published)
		VALUES (gen_random_uuid(), $1::text::uuid, 1, false) RETURNING id::text`, draft.ID).Scan(&draft.VersionID)
	if err != nil {
		return domain.PipelineDraft{}, err
	}
	if _, err := tx.Exec(ctx, `
		UPDATE job_postings SET pipeline_id = $1::text::uuid WHERE id = $2::text::uuid`, draft.ID, posting); err != nil {
		return domain.PipelineDraft{}, err
	}
	return draft, tx.Commit(ctx)
}

func (r *PipelineRepository) Stages(ctx context.Context, company, version string) ([]domain.Stage, error) {
	// Keep ownership and stage reads in one snapshot, including for an empty draft.
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if _, err := pipelineVersion(ctx, tx, company, version, false); err != nil {
		return nil, err
	}
	rows, err := tx.Query(ctx, `
		SELECT s.id::text, s.name, s.module_type, s.scheduled_start_at, s.scheduled_end_at,
			s.interview_configuration,
			ARRAY(SELECT tc.test_module_id::text FROM test_configurations tc WHERE tc.stage_id = s.id
				ORDER BY (tc.configuration->>'position')::integer)
		FROM pipeline_stages s WHERE s.version_id = $1::text::uuid
		ORDER BY s.stage_order`, version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	stages := make([]domain.Stage, 0)
	for rows.Next() {
		var stage domain.Stage
		err := rows.Scan(&stage.ID, &stage.Name, &stage.Type, &stage.OpensAt, &stage.Deadline,
			&stage.Interview, &stage.ModuleIDs)
		if err != nil {
			return nil, err
		}
		stages = append(stages, stage)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return stages, tx.Commit(ctx)
}

func (r *PipelineRepository) ReplaceStages(ctx context.Context, company, version string, stages []domain.Stage) ([]domain.Stage, error) {
	if err := domain.ValidateStages(stages); err != nil {
		return nil, err
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	published, err := pipelineVersion(ctx, tx, company, version, true)
	if err != nil {
		return nil, err
	}
	if published {
		return nil, domain.ErrPipelineLocked
	}

	definitions := make(map[string]domain.Definition)
	for _, stage := range stages {
		if stage.Interview != nil {
			if err := validateInterviewers(ctx, tx, company, stage.Interview.InterviewerIDs); err != nil {
				return nil, err
			}
		}
		for _, id := range stage.ModuleIDs {
			if _, exists := definitions[id]; exists {
				continue
			}
			definition, err := moduleSnapshot(ctx, tx, company, id)
			if err != nil {
				return nil, err
			}
			definitions[id] = definition
		}
	}

	if _, err := tx.Exec(ctx, `
		DELETE FROM test_configurations WHERE stage_id IN
			(SELECT id FROM pipeline_stages WHERE version_id = $1::text::uuid)`, version); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM pipeline_stages WHERE version_id = $1::text::uuid`, version); err != nil {
		return nil, err
	}
	for i := range stages {
		if err := insertStage(ctx, tx, version, i+1, &stages[i], definitions); err != nil {
			return nil, err
		}
	}
	return stages, tx.Commit(ctx)
}

func pipelineVersion(ctx context.Context, tx pgx.Tx, company, version string, lock bool) (bool, error) {
	query := `
		SELECT COALESCE(v.is_published, false) OR v.published_at IS NOT NULL OR p.published_version_id IS NOT DISTINCT FROM v.id
		FROM pipeline_versions v JOIN pipelines p ON p.id = v.pipeline_id
		JOIN job_postings j ON j.id = p.posting_id
		WHERE v.id = $1::text::uuid AND j.company_id = $2::text::uuid`
	if lock {
		query += " FOR UPDATE OF v"
	}
	var published bool
	if err := tx.QueryRow(ctx, query, version, company).Scan(&published); err != nil {
		return false, pipelineError(err)
	}

	var legacy bool
	err := tx.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM pipeline_stages WHERE version_id = $1::text::uuid AND stage_order IS NULL)
			OR EXISTS (SELECT 1 FROM stage_transitions WHERE pipeline_version_id = $1::text::uuid)`, version).Scan(&legacy)
	if err == nil && legacy {
		return false, domain.ErrPipelineLocked
	}
	return published, err
}

func moduleSnapshot(ctx context.Context, tx pgx.Tx, company, id string) (domain.Definition, error) {
	module, err := scanModule(tx.QueryRow(ctx, `
		SELECT `+moduleColumns+` FROM modules WHERE id = $1::text::uuid
			AND company_id = $2::text::uuid AND authoring_definition IS NOT NULL FOR SHARE`, id, company))
	if errors.Is(err, domain.ErrModuleNotFound) {
		return domain.Definition{}, &domain.InvalidPipelineError{Reason: fmt.Sprintf("module %s is unavailable to this company", id)}
	}
	if err != nil {
		return domain.Definition{}, err
	}
	if module.Published || module.CopiedFrom != nil {
		return module.Definition, nil
	}
	if err := prepareModuleDefinition(ctx, tx, &module.Definition); err != nil {
		var invalid *domain.InvalidModuleError
		if errors.As(err, &invalid) {
			return domain.Definition{}, &domain.InvalidPipelineError{Reason: fmt.Sprintf("module %s: %s", id, invalid.Reason)}
		}
		return domain.Definition{}, err
	}
	return module.Definition, nil
}

func validateInterviewers(ctx context.Context, tx pgx.Tx, company string, ids []string) error {
	var count int
	err := tx.QueryRow(ctx, `
		SELECT count(*) FROM interviewers WHERE company_id = $1::text::uuid
			AND id::text = ANY($2::text[]) AND deleted_at IS NULL AND deactivated_at IS NULL`, company, ids).Scan(&count)
	if err == nil && count != len(ids) {
		return &domain.InvalidPipelineError{Reason: "every interviewer must be active and belong to this company"}
	}
	return err
}

func insertStage(ctx context.Context, tx pgx.Tx, version string, position int, stage *domain.Stage, definitions map[string]domain.Definition) error {
	rule := "completion"
	var threshold *int
	if stage.Interview != nil {
		rule = "manual"
		if !stage.Interview.Behavioral {
			rule = "threshold"
			threshold = stage.Interview.Threshold
		}
	}
	err := tx.QueryRow(ctx, `
		INSERT INTO pipeline_stages (id, version_id, name, module_type, stage_order,
			scheduled_start_at, scheduled_end_at, interview_configuration, advancement_rule,
			score_threshold, below_threshold_action)
		VALUES (gen_random_uuid(), $1::text::uuid, $2, $3, $4, $5, $6, $7, $8, $9, 'reject')
		RETURNING id::text`, version, stage.Name, stage.Type, position,
		stage.OpensAt, stage.Deadline, stage.Interview, rule, threshold).Scan(&stage.ID)
	if err != nil {
		return err
	}
	for i, id := range stage.ModuleIDs {
		_, err := tx.Exec(ctx, `
			INSERT INTO test_configurations (id, stage_id, test_module_id, active_time_limit, configuration)
			VALUES (gen_random_uuid(), $1::text::uuid, $2::text::uuid, make_interval(mins => $3),
				jsonb_build_object('position', $4::integer, 'definition', $5::jsonb))`,
			stage.ID, id, definitions[id].ActiveMinutes, i, definitions[id])
		if err != nil {
			return err
		}
	}
	return nil
}

func pipelineError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrPipelineNotFound
	}
	return err
}

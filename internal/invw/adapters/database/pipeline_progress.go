package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CSCE-3701-Amara/amara/internal/invw/domain"
)

type CandidatePipelineProgressRepository struct {
	pool *pgxpool.Pool
}

func NewCandidatePipelineProgressRepository(pool *pgxpool.Pool) *CandidatePipelineProgressRepository {
	return &CandidatePipelineProgressRepository{pool: pool}
}

func (r *CandidatePipelineProgressRepository) GetByCandidateApplication(ctx context.Context, candidateID, applicationID string) (domain.CandidatePipelineProgress, error) {
	if r == nil || r.pool == nil {
		return domain.CandidatePipelineProgress{}, fmt.Errorf("candidate pipeline progress database is not configured")
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return domain.CandidatePipelineProgress{}, fmt.Errorf("begin candidate pipeline progress read: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	result := domain.CandidatePipelineProgress{ApplicationID: applicationID}
	var pipelineExecutionID string
	err = tx.QueryRow(ctx, `
		SELECT pv.id::text, pv.version_number, pv.published_at, execution.id::text
		FROM applications a
		JOIN pipeline_executions execution ON execution.application_id = a.id
		JOIN pipeline_versions pv ON pv.id = execution.pipeline_version_id
		WHERE a.id = $1 AND a.cand_id = $2 AND pv.published_at IS NOT NULL`,
		applicationID, candidateID,
	).Scan(&result.Pipeline.ID, &result.Pipeline.Version, &result.Pipeline.PublishedAt, &pipelineExecutionID)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.CandidatePipelineProgress{}, fmt.Errorf("%w: no published pipeline is assigned to this application", domain.ErrPipelineProgressNotFound)
	}
	if err != nil {
		return domain.CandidatePipelineProgress{}, fmt.Errorf("load assigned published pipeline: %w", err)
	}

	rows, err := tx.Query(ctx, `
		SELECT stage.id::text, stage.name, stage.module_type,
		       execution.status, execution.completed_at
		FROM pipeline_stages stage
		LEFT JOIN stage_executions execution
		  ON execution.stage_id = stage.id AND execution.pipeline_execution_id = $2
		WHERE stage.version_id = $1
		ORDER BY stage.position, stage.id`, result.Pipeline.ID, pipelineExecutionID)
	if err != nil {
		return domain.CandidatePipelineProgress{}, fmt.Errorf("query published pipeline stages: %w", err)
	}
	for rows.Next() {
		var stage domain.PipelineStage
		var executionStatus *string
		var completedAt *time.Time
		if err := rows.Scan(&stage.ID, &stage.Name, &stage.Type, &executionStatus, &completedAt); err != nil {
			rows.Close()
			return domain.CandidatePipelineProgress{}, fmt.Errorf("scan published pipeline stage: %w", err)
		}
		if executionStatus != nil {
			completed := *executionStatus == "passed" || *executionStatus == "rejected" ||
				*executionStatus == "cancelled" || *executionStatus == "missed"
			active := *executionStatus == "available" || *executionStatus == "in_progress" ||
				*executionStatus == "pending_grading" || *executionStatus == "held"
			stage.Progress = &domain.StageProgress{Completed: completed, Active: active}
			if completed {
				stage.Progress.Passed = *executionStatus == "passed"
				stage.Progress.CompletedAt = completedAt
			}
		}
		stage.Awards = make([]domain.AmaraAward, 0)
		result.Stages = append(result.Stages, stage)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return domain.CandidatePipelineProgress{}, fmt.Errorf("read published pipeline stages: %w", err)
	}
	rows.Close()

	rows, err = tx.Query(ctx, `
		SELECT link.from_stage_id::text, link.to_stage_id::text
		FROM stage_transitions link
		WHERE link.pipeline_version_id = $1
		ORDER BY link.from_stage_id, link.to_stage_id`, result.Pipeline.ID)
	if err != nil {
		return domain.CandidatePipelineProgress{}, fmt.Errorf("query published pipeline links: %w", err)
	}
	result.Links = make([]domain.PipelineLink, 0)
	for rows.Next() {
		var link domain.PipelineLink
		if err := rows.Scan(&link.FromStageID, &link.ToStageID); err != nil {
			rows.Close()
			return domain.CandidatePipelineProgress{}, fmt.Errorf("scan published pipeline link: %w", err)
		}
		result.Links = append(result.Links, link)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return domain.CandidatePipelineProgress{}, fmt.Errorf("read published pipeline links: %w", err)
	}
	rows.Close()

	rows, err = tx.Query(ctx, `
		SELECT award.id::text, execution.stage_id::text,
		       COALESCE(award.company_id::text, ''),
		       COALESCE(NULLIF(BTRIM(company.display_name), ''), company.legal_name, ''),
		       award.amara_title,
		       COALESCE(award.definition_snapshot->>'crest', ''),
		       COALESCE(award.definition_snapshot->>'domain', ''),
		       award.level, award.definition_snapshot->>'level_description',
		       award.strengths, award.areas_for_development,
		       CASE WHEN jsonb_typeof(award.pass_decision->'evidence') = 'array'
		            THEN award.pass_decision->'evidence' ELSE '[]'::jsonb END,
		       award.awarded_at
		FROM candidate_amarat award
		JOIN stage_executions execution
		  ON execution.id = award.source_stage_execution_id
		 AND execution.pipeline_execution_id = $3
		 AND execution.status = 'passed'
		LEFT JOIN companies company ON company.id = award.company_id
		WHERE award.application_id = $1 AND award.cand_id = $2
		ORDER BY award.awarded_at, award.id`, applicationID, candidateID, pipelineExecutionID)
	if err != nil {
		return domain.CandidatePipelineProgress{}, fmt.Errorf("query candidate Amara awards: %w", err)
	}
	defer rows.Close()
	stages := make(map[string]*domain.PipelineStage, len(result.Stages))
	for i := range result.Stages {
		stages[result.Stages[i].ID] = &result.Stages[i]
	}
	for rows.Next() {
		var award domain.AmaraAward
		var stageID string
		var strengths, developmentAreas *string
		var evidence []byte
		if err := rows.Scan(
			&award.ID, &stageID, &award.CompanyID, &award.CompanyName, &award.Name,
			&award.Crest, &award.SkillDomain, &award.Level, &award.LevelDescription,
			&strengths, &developmentAreas, &evidence, &award.AwardedAt,
		); err != nil {
			return domain.CandidatePipelineProgress{}, fmt.Errorf("scan candidate Amara award: %w", err)
		}
		if err := json.Unmarshal(evidence, &award.Evidence); err != nil {
			return domain.CandidatePipelineProgress{}, fmt.Errorf("decode Amara evidence: %w", err)
		}
		award.Strengths = awardTextList(strengths)
		award.DevelopmentAreas = awardTextList(developmentAreas)
		if award.Evidence == nil {
			award.Evidence = []string{}
		}
		if stage := stages[stageID]; stage != nil {
			stage.Awards = append(stage.Awards, award)
		}
	}
	if err := rows.Err(); err != nil {
		return domain.CandidatePipelineProgress{}, fmt.Errorf("read candidate Amara awards: %w", err)
	}
	return result, nil
}

func awardTextList(value *string) []string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return []string{}
	}
	return []string{*value}
}

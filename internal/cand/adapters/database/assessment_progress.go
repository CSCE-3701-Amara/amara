package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CSCE-3701-Amara/amara/internal/cand/domain"
)

type AssessmentProgressRepository struct {
	pool *pgxpool.Pool
}

func NewAssessmentProgressRepository(pool *pgxpool.Pool) *AssessmentProgressRepository {
	return &AssessmentProgressRepository{pool: pool}
}

func (r *AssessmentProgressRepository) ListByCandidateApplication(ctx context.Context, candidateID, applicationID string) ([]domain.AssessmentProgress, bool, error) {
	if r == nil || r.pool == nil {
		return nil, false, fmt.Errorf("candidate assessment progress database is not configured")
	}
	var applicationFound bool
	if err := r.pool.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM applications WHERE id = $1 AND cand_id = $2)`,
		applicationID, candidateID,
	).Scan(&applicationFound); err != nil {
		return nil, false, fmt.Errorf("check candidate application ownership: %w", err)
	}
	if !applicationFound {
		return nil, false, nil
	}

	rows, err := r.pool.Query(ctx, `
		SELECT aa.id::text, aa.status, aa.supports_interruption, aa.saved_progress, aa.updated_at
		FROM application_assessments aa
		JOIN applications a ON a.id = aa.application_id
		WHERE a.id = $1 AND a.cand_id = $2
		ORDER BY aa.created_at, aa.id`, applicationID, candidateID)
	if err != nil {
		return nil, false, fmt.Errorf("query candidate assessment progress: %w", err)
	}
	defer rows.Close()

	assessments := make([]domain.AssessmentProgress, 0)
	for rows.Next() {
		var assessment domain.AssessmentProgress
		if err := rows.Scan(&assessment.ID, &assessment.Status, &assessment.SupportsInterruption, &assessment.SavedProgress, &assessment.UpdatedAt); err != nil {
			return nil, false, fmt.Errorf("scan candidate assessment progress: %w", err)
		}
		assessments = append(assessments, assessment)
	}
	if err := rows.Err(); err != nil {
		return nil, false, fmt.Errorf("read candidate assessment progress: %w", err)
	}
	return assessments, true, nil
}

func (r *AssessmentProgressRepository) SaveProgress(ctx context.Context, candidateID, applicationID, assessmentID string, progress json.RawMessage) (domain.AssessmentProgress, error) {
	if r == nil || r.pool == nil {
		return domain.AssessmentProgress{}, fmt.Errorf("candidate assessment progress database is not configured")
	}
	var assessment domain.AssessmentProgress
	err := r.pool.QueryRow(ctx, `
		UPDATE application_assessments aa
		SET saved_progress = $4::jsonb,
		    status = CASE WHEN aa.status = 'not_started' THEN 'in_progress' ELSE aa.status END,
		    updated_at = now()
		FROM applications a
		WHERE aa.application_id = a.id
		  AND a.id = $1 AND a.cand_id = $2 AND aa.id = $3
		  AND aa.supports_interruption
		  AND aa.status IN ('not_started', 'in_progress')
		RETURNING aa.id::text, aa.status, aa.supports_interruption, aa.saved_progress, aa.updated_at`,
		applicationID, candidateID, assessmentID, string(progress),
	).Scan(&assessment.ID, &assessment.Status, &assessment.SupportsInterruption, &assessment.SavedProgress, &assessment.UpdatedAt)
	if err == nil {
		return assessment, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.AssessmentProgress{}, fmt.Errorf("update candidate assessment progress: %w", err)
	}

	var supportsInterruption bool
	var status domain.AssessmentStatus
	err = r.pool.QueryRow(ctx, `
		SELECT aa.status, aa.supports_interruption
		FROM application_assessments aa
		JOIN applications a ON a.id = aa.application_id
		WHERE a.id = $1 AND a.cand_id = $2 AND aa.id = $3`,
		applicationID, candidateID, assessmentID,
	).Scan(&status, &supportsInterruption)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.AssessmentProgress{}, domain.ErrAssessmentNotFound
	}
	if err != nil {
		return domain.AssessmentProgress{}, fmt.Errorf("check candidate assessment state: %w", err)
	}
	if !supportsInterruption || (status != domain.AssessmentStatusNotStarted && status != domain.AssessmentStatusInProgress) {
		return domain.AssessmentProgress{}, domain.ErrAssessmentNotResumable
	}
	return domain.AssessmentProgress{}, fmt.Errorf("assessment progress changed while saving")
}

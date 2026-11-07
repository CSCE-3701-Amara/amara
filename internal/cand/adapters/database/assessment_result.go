package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CSCE-3701-Amara/amara/internal/cand/domain"
)

type AssessmentResultRepository struct {
	pool *pgxpool.Pool
}

func NewAssessmentResultRepository(pool *pgxpool.Pool) *AssessmentResultRepository {
	return &AssessmentResultRepository{pool: pool}
}

func (r *AssessmentResultRepository) ListByCandidateApplication(ctx context.Context, candidateID, applicationID string) ([]domain.AssessmentResult, bool, error) {
	if r == nil || r.pool == nil {
		return nil, false, fmt.Errorf("candidate assessment result database is not configured")
	}
	rows, err := r.pool.Query(ctx, `
		SELECT i.id::text,
		       CASE WHEN ar.id IS NOT NULL OR t.evaluated_at IS NOT NULL THEN 'evaluated' ELSE 'pending' END,
		       CASE WHEN ar.disclosed_at IS NOT NULL OR ar.is_disclosed IS TRUE
		            THEN ar.passed ELSE NULL END,
		       CASE WHEN ar.disclosed_at IS NOT NULL OR ar.is_disclosed IS TRUE
		            THEN coalesce(ar.scored_at, t.evaluated_at) ELSE NULL END,
		       CASE WHEN ar.feedback_approved_at IS NOT NULL
		                  AND (ar.disclosed_at IS NOT NULL OR ar.is_disclosed IS TRUE)
		            THEN ar.feedback ELSE NULL END,
		       '[]'::jsonb,
		       '[]'::jsonb,
		       CASE WHEN ar.feedback_approved_at IS NOT NULL
		                  AND (ar.disclosed_at IS NOT NULL OR ar.is_disclosed IS TRUE)
		            THEN ar.feedback_is_ai_generated ELSE NULL END
		FROM applications a
		JOIN assessment_instances i ON i.application_id = a.id
		LEFT JOIN assessment_attempts t ON t.instance_id = i.id
		LEFT JOIN assessment_results ar ON ar.instance_id = i.id
		WHERE a.id = $1 AND a.cand_id = $2
		  AND (t.submitted_at IS NOT NULL OR ar.id IS NOT NULL)
		ORDER BY coalesce(t.submitted_at, i.deadline), i.id`, applicationID, candidateID)
	if err != nil {
		return nil, false, fmt.Errorf("query candidate assessment results: %w", err)
	}
	defer rows.Close()

	results := make([]domain.AssessmentResult, 0)
	applicationFound := false
	for rows.Next() {
		applicationFound = true
		var result domain.AssessmentResult
		var passed *bool
		var evaluatedAt *time.Time
		var feedbackText *string
		var strengthsJSON, improvementsJSON []byte
		var aiGenerated *bool
		if err := rows.Scan(
			&result.AssessmentID, &result.Status, &passed, &evaluatedAt,
			&feedbackText, &strengthsJSON, &improvementsJSON, &aiGenerated,
		); err != nil {
			return nil, false, fmt.Errorf("scan candidate assessment result: %w", err)
		}
		if passed != nil {
			result.Outcome = &domain.AssessmentOutcome{Passed: *passed}
			result.EvaluatedAt = evaluatedAt
		}
		if feedbackText != nil || strengthsJSON != nil || improvementsJSON != nil || aiGenerated != nil {
			feedback := &domain.AssessmentFeedback{Text: feedbackText}
			if len(strengthsJSON) > 0 {
				if err := json.Unmarshal(strengthsJSON, &feedback.Strengths); err != nil {
					return nil, false, fmt.Errorf("decode assessment feedback strengths: %w", err)
				}
			}
			if len(improvementsJSON) > 0 {
				if err := json.Unmarshal(improvementsJSON, &feedback.AreasForImprovement); err != nil {
					return nil, false, fmt.Errorf("decode assessment improvement areas: %w", err)
				}
			}
			if aiGenerated != nil {
				feedback.AIGenerated = *aiGenerated
			}
			result.Feedback = feedback
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, false, fmt.Errorf("read candidate assessment results: %w", err)
	}
	if applicationFound {
		return results, true, nil
	}
	if err := r.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM applications WHERE id = $1 AND cand_id = $2)`, applicationID, candidateID).Scan(&applicationFound); err != nil {
		return nil, false, fmt.Errorf("check candidate application ownership: %w", err)
	}
	return results, applicationFound, nil
}

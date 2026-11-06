package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CSCE-3701-Amara/amara/internal/cand/domain"
)

type ApplicationTrackingRepository struct {
	pool *pgxpool.Pool
}

func NewApplicationTrackingRepository(pool *pgxpool.Pool) *ApplicationTrackingRepository {
	return &ApplicationTrackingRepository{pool: pool}
}

const trackingColumns = `
	a.id::text, a.reference_number, p.id::text, p.title,
	COALESCE(NULLIF(c.display_name, ''), c.legal_name), a.submitted_at,
	CASE
		WHEN a.withdrawn_at IS NOT NULL THEN 'withdrawn'
		WHEN a.screened_out IS TRUE THEN 'screened_out'
		WHEN e.status <> 'active' THEN e.status
		WHEN a.final_decision IS FALSE THEN 'rejected'
		WHEN a.final_decision IS TRUE THEN 'completed'
		WHEN e.status = 'active' THEN 'in_progress'
		ELSE COALESCE(NULLIF(a.status, ''), 'submitted')
	END AS status`

const trackingFrom = `
	FROM applications a
	JOIN job_postings p ON p.id = a.posting_id
	JOIN companies c ON c.id = p.company_id
	LEFT JOIN pipeline_executions e ON e.application_id = a.id
	WHERE a.cand_id = $1::uuid`

func (r *ApplicationTrackingRepository) List(ctx context.Context, candidateID string) ([]domain.ApplicationSummary, error) {
	rows, err := r.pool.Query(ctx, "SELECT "+trackingColumns+trackingFrom+
		" ORDER BY COALESCE(a.submitted_at, a.created_at) DESC NULLS LAST, a.id DESC", candidateID)
	if err != nil {
		return nil, fmt.Errorf("query application history: %w", err)
	}
	defer rows.Close()

	applications := make([]domain.ApplicationSummary, 0)
	for rows.Next() {
		var summary domain.ApplicationSummary
		if err := rows.Scan(summaryFields(&summary)...); err != nil {
			return nil, fmt.Errorf("scan application history: %w", err)
		}
		summary.Closed = applicationClosed(summary.Status)
		applications = append(applications, summary)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read application history: %w", err)
	}
	return applications, nil
}

func (r *ApplicationTrackingRepository) Get(ctx context.Context, candidateID, applicationID string) (domain.ApplicationDetail, error) {
	var detail domain.ApplicationDetail
	fields := append(summaryFields(&detail.ApplicationSummary),
		&detail.PipelineExecutionID, &detail.FinalDecision, &detail.FinalFeedback, &detail.RejectionFeedback)
	err := r.pool.QueryRow(ctx, "SELECT "+trackingColumns+`,
		e.id::text, a.final_decision, a.final_feedback,
		CASE WHEN e.feedback_released_at <= now() THEN e.rejection_feedback END`+
		trackingFrom+" AND a.id = $2::uuid", candidateID, applicationID).Scan(fields...)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ApplicationDetail{}, domain.ErrApplicationNotFound
	}
	if err != nil {
		return domain.ApplicationDetail{}, fmt.Errorf("read application detail: %w", err)
	}

	detail.Closed = applicationClosed(detail.Status)
	if !detail.Closed {
		detail.FinalFeedback = nil
		detail.RejectionFeedback = nil
	}
	return detail, nil
}

func summaryFields(summary *domain.ApplicationSummary) []any {
	return []any{&summary.ID, &summary.ReferenceNumber, &summary.PostingID, &summary.JobTitle,
		&summary.CompanyName, &summary.SubmittedAt, &summary.Status}
}

func applicationClosed(status string) bool {
	switch status {
	case "screened_out", "completed", "rejected", "withdrawn", "cancelled", "closed":
		return true
	default:
		return false
	}
}

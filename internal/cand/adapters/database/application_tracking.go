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
	a.id::text AS id, a.reference_number, p.id::text AS posting_id, p.title,
	COALESCE(NULLIF(c.display_name, ''), c.legal_name) AS company_name,
	a.created_at, a.updated_at, a.submitted_at, a.withdrawn_at, a.screened_out, a.screening_completed_at,
	COALESCE(w.anonymity_enabled, FALSE) AND NOT EXISTS (SELECT 1 FROM identity_reveal_records r WHERE r.application_id = a.id) AS identity_hidden,
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
	LEFT JOIN pipeline_versions v ON v.id = e.pipeline_version_id
	LEFT JOIN pipelines w ON w.id = v.pipeline_id
	WHERE a.cand_id = $1::uuid`

func (r *ApplicationTrackingRepository) List(ctx context.Context, candidateID string, filter domain.ApplicationFilter) (domain.ApplicationPage, error) {
	page := domain.ApplicationPage{
		Data:       make([]domain.ApplicationSummary, 0),
		Pagination: domain.ApplicationPagination{Limit: filter.Limit, Offset: filter.Offset},
	}
	// Count and read the page against the same snapshot, including when the offset is past the end.
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return domain.ApplicationPage{}, fmt.Errorf("begin application history read: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const filtered = " FROM (SELECT " + trackingColumns + trackingFrom +
		") AS history WHERE ($2::text = '' OR status = $2)"
	err = tx.QueryRow(ctx, "SELECT count(*)"+filtered, candidateID, filter.Status).Scan(&page.Pagination.Total)
	if err != nil {
		return domain.ApplicationPage{}, fmt.Errorf("count application history: %w", err)
	}

	rows, err := tx.Query(ctx, "SELECT *"+filtered+
		" ORDER BY COALESCE(submitted_at, created_at) DESC NULLS LAST, id DESC LIMIT $3 OFFSET $4",
		candidateID, filter.Status, filter.Limit, filter.Offset)
	if err != nil {
		return domain.ApplicationPage{}, fmt.Errorf("query application history: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var summary domain.ApplicationSummary
		if err := rows.Scan(summaryFields(&summary)...); err != nil {
			return domain.ApplicationPage{}, fmt.Errorf("scan application history: %w", err)
		}
		summary.Closed = applicationClosed(summary.Status)
		page.Data = append(page.Data, summary)
	}
	if err := rows.Err(); err != nil {
		return domain.ApplicationPage{}, fmt.Errorf("read application history: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.ApplicationPage{}, fmt.Errorf("complete application history read: %w", err)
	}
	return page, nil
}

func (r *ApplicationTrackingRepository) Get(ctx context.Context, candidateID, applicationID string) (domain.ApplicationDetail, error) {
	var detail domain.ApplicationDetail
	fields := append(summaryFields(&detail.ApplicationSummary),
		&detail.PipelineExecutionID, &detail.FinalDecision, &detail.FinalFeedback, &detail.RejectionFeedback)
	err := r.pool.QueryRow(ctx, "SELECT "+trackingColumns+`,
		e.id::text, a.final_decision,
		CASE WHEN e.feedback_released_at <= now() THEN a.final_feedback END,
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
		&summary.CompanyName, &summary.CreatedAt, &summary.UpdatedAt, &summary.SubmittedAt,
		&summary.WithdrawnAt, &summary.ScreenedOut, &summary.ScreeningCompletedAt, &summary.IdentityHidden, &summary.Status}
}

func applicationClosed(status string) bool {
	switch status {
	case "screened_out", "completed", "rejected", "withdrawn", "cancelled", "closed":
		return true
	default:
		return false
	}
}

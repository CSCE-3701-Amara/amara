package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostingLifecycleRepository persists JOB-owned posting state transitions.
type PostingLifecycleRepository struct {
	pool *pgxpool.Pool
}

func NewPostingLifecycleRepository(pool *pgxpool.Pool) *PostingLifecycleRepository {
	return &PostingLifecycleRepository{pool: pool}
}

func (r *PostingLifecycleRepository) ArchiveExpired(ctx context.Context) error {
	if r == nil || r.pool == nil {
		return fmt.Errorf("job posting database is not configured")
	}
	_, err := r.pool.Exec(ctx, `UPDATE job_postings SET status = 'archived', archived_at = COALESCE(archived_at, now()), updated_at = now() WHERE status = 'published' AND deleted_at IS NULL AND application_deadline IS NOT NULL AND application_deadline <= now()`)
	if err != nil {
		return fmt.Errorf("update expired job postings: %w", err)
	}
	return nil
}

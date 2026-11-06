package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/job/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ApplicationRepository persists submissions and their initial screening result.
type ApplicationRepository struct{ pool *pgxpool.Pool }

func NewApplicationRepository(pool *pgxpool.Pool) *ApplicationRepository {
	return &ApplicationRepository{pool: pool}
}

func (r *ApplicationRepository) Apply(ctx context.Context, postingID, candidateID string) (domain.Application, error) {
	if r == nil || r.pool == nil {
		return domain.Application{}, errors.New("application database is not configured")
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Application{}, fmt.Errorf("begin application submission: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	// Serialize per posting so simultaneous requests cannot exceed max_applications.
	var postingStatus string
	var opening, deadline, deletedAt *time.Time
	var maxApplications *int32
	err = tx.QueryRow(ctx, `SELECT status,application_opening,application_deadline,max_applications,deleted_at FROM job_postings WHERE id=$1::uuid FOR UPDATE`, postingID).Scan(&postingStatus, &opening, &deadline, &maxApplications, &deletedAt)
	if errors.Is(err, pgx.ErrNoRows) || deletedAt != nil {
		return domain.Application{}, domain.ErrPostingNotFound
	}
	if err != nil {
		return domain.Application{}, fmt.Errorf("load posting: %w", err)
	}
	var now time.Time
	if err := tx.QueryRow(ctx, "SELECT clock_timestamp()").Scan(&now); err != nil {
		return domain.Application{}, fmt.Errorf("read database time: %w", err)
	}
	if !domain.ApplicationIsOpen(postingStatus, opening, deadline, now) {
		return domain.Application{}, domain.ErrApplicationPostingUnavailable
	}
	var eligible bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM candidates WHERE id=$1::uuid AND deleted_at IS NULL AND deactivated_at IS NULL AND COALESCE(banned,FALSE)=FALSE)`, candidateID).Scan(&eligible); err != nil {
		return domain.Application{}, fmt.Errorf("check candidate eligibility: %w", err)
	}
	if !eligible {
		return domain.Application{}, domain.ErrCandidateUnavailable
	}
	var duplicate bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM applications WHERE posting_id=$1::uuid AND cand_id=$2::uuid)`, postingID, candidateID).Scan(&duplicate); err != nil {
		return domain.Application{}, fmt.Errorf("check duplicate application: %w", err)
	}
	if duplicate {
		return domain.Application{}, domain.ErrDuplicateApplication
	}
	if maxApplications != nil {
		var count int64
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM applications WHERE posting_id=$1::uuid`, postingID).Scan(&count); err != nil {
			return domain.Application{}, fmt.Errorf("count applications: %w", err)
		}
		if domain.ApplicationCapReached(maxApplications, count) {
			return domain.Application{}, domain.ErrApplicationLimitReached
		}
	}
	// Exact case-insensitive matching uses current candidate_amarat. Free-text posting
	// requirements are informational and are not evaluated. Immutable resume/profile
	// snapshots are not implemented; failed criteria are retained.
	requiredAmarat, err := applicationAmaraTitles(ctx, tx, `SELECT amara_title FROM job_posting_amarat WHERE posting_id=$1::uuid ORDER BY lower(amara_title),amara_title`, postingID)
	if err != nil {
		return domain.Application{}, fmt.Errorf("load required Amarat: %w", err)
	}
	candidateAmarat, err := applicationAmaraTitles(ctx, tx, `SELECT amara_title FROM candidate_amarat WHERE cand_id=$1::uuid`, candidateID)
	if err != nil {
		return domain.Application{}, fmt.Errorf("load candidate Amarat: %w", err)
	}
	failed := domain.MissingAmarat(requiredAmarat, candidateAmarat)

	appStatus := domain.ApplicationAdmitted
	screenedOut := len(failed) > 0
	if screenedOut {
		appStatus = domain.ApplicationScreenedOut
	}
	var result domain.Application
	err = tx.QueryRow(ctx, `INSERT INTO applications(id,posting_id,cand_id,status,created_at,updated_at,submitted_at,screened_out,screening_completed_at,failed_amarat)
 VALUES(gen_random_uuid(),$1::uuid,$2::uuid,$3,clock_timestamp(),clock_timestamp(),clock_timestamp(),$4,clock_timestamp(),$5)
 RETURNING id::text,reference_number,posting_id::text,status,submitted_at,screened_out,screening_completed_at,failed_amarat`,
		postingID, candidateID, appStatus, screenedOut, failed).Scan(&result.ID, &result.ReferenceNumber, &result.PostingID, &result.Status, &result.SubmittedAt, &result.ScreenedOut, &result.ScreeningCompletedAt, &result.FailedAmarat)
	if err != nil {
		return domain.Application{}, fmt.Errorf("save screened application: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Application{}, fmt.Errorf("commit submission: %w", err)
	}
	// TODO(INVW): send admitted applications into the published pipeline. This checkout
	// has no JOB-to-INVW admission event/outbox consumer; JOB admission is recorded only.
	return result, nil
}

func applicationAmaraTitles(ctx context.Context, tx pgx.Tx, query string, id string) ([]string, error) {
	rows, err := tx.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	titles := make([]string, 0)
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, err
		}
		titles = append(titles, title)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return titles, nil
}

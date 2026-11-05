package database

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CSCE-3701-Amara/amara/internal/job/domain"
)

// PostingManagementRepository persists JOB-owned posting configuration and lifecycle.
type PostingManagementRepository struct{ pool *pgxpool.Pool }

func NewPostingManagementRepository(pool *pgxpool.Pool) *PostingManagementRepository {
	return &PostingManagementRepository{pool: pool}
}

const postingColumns = `id::text, company_id::text, title, description, language, currency,
	seniority_level, employment_type, location, salary_min, salary_max, salary_visible,
	open_positions, application_opening, application_deadline, max_applications, status,
	last_autosaved_at, published_at, archived_at, created_at, updated_at`

type rowQuerier interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

func scanPosting(row pgx.Row) (domain.Posting, error) {
	var p domain.Posting
	var seniority, employmentType, status pgtype.Text
	var createdAt, updatedAt pgtype.Timestamp
	err := row.Scan(&p.ID, &p.CompanyID, &p.Title, &p.Description, &p.Language, &p.Currency,
		&seniority, &employmentType, &p.Location, &p.SalaryMin, &p.SalaryMax,
		&p.SalaryVisible, &p.OpenPositions, &p.ApplicationOpening, &p.ApplicationDeadline,
		&p.MaxApplications, &status, &p.LastAutosavedAt, &p.PublishedAt, &p.ArchivedAt,
		&createdAt, &updatedAt)
	if seniority.Valid {
		p.SeniorityLevel = seniority.String
	}
	if employmentType.Valid {
		p.EmploymentType = employmentType.String
	}
	if status.Valid {
		p.Status = status.String
	}
	if createdAt.Valid {
		p.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		p.UpdatedAt = updatedAt.Time
	}
	return p, err
}

func (r *PostingManagementRepository) Create(ctx context.Context, companyID string, p domain.Posting) (domain.Posting, error) {
	if r == nil || r.pool == nil {
		return domain.Posting{}, errors.New("job posting database is not configured")
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Posting{}, fmt.Errorf("begin posting creation: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	row, err := scanPosting(tx.QueryRow(ctx, `INSERT INTO job_postings
		(id, company_id, title, description, language, currency, seniority_level, employment_type,
		 location, salary_min, salary_max, salary_visible, open_positions, application_opening,
		 application_deadline, max_applications, status, created_at, updated_at)
		VALUES (gen_random_uuid(), $1::uuid, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, 'draft', now(), now())
		RETURNING `+postingColumns,
		companyID, strings.TrimSpace(p.Title), p.Description, p.Language, p.Currency,
		strings.TrimSpace(p.SeniorityLevel), p.EmploymentType, p.Location, p.SalaryMin,
		p.SalaryMax, p.SalaryVisible, p.OpenPositions, p.ApplicationOpening,
		p.ApplicationDeadline, p.MaxApplications))
	if err != nil {
		return domain.Posting{}, fmt.Errorf("insert job posting: %w", err)
	}
	if err := replaceScreeningData(ctx, tx, row.ID, p); err != nil {
		return domain.Posting{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Posting{}, fmt.Errorf("commit posting creation: %w", err)
	}
	return r.Get(ctx, companyID, row.ID)
}

func (r *PostingManagementRepository) Get(ctx context.Context, companyID, id string) (domain.Posting, error) {
	if r == nil || r.pool == nil {
		return domain.Posting{}, errors.New("job posting database is not configured")
	}
	p, err := scanPosting(r.pool.QueryRow(ctx, `SELECT `+postingColumns+
		` FROM job_postings WHERE id = $1::uuid AND company_id = $2::uuid AND deleted_at IS NULL`, id, companyID))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Posting{}, domain.ErrPostingNotFound
	}
	if err != nil {
		return domain.Posting{}, fmt.Errorf("get job posting: %w", err)
	}
	if err := loadScreeningData(ctx, r.pool, &p); err != nil {
		return domain.Posting{}, err
	}
	return p, nil
}

func (r *PostingManagementRepository) ListByCompany(ctx context.Context, companyID string) ([]domain.Posting, error) {
	if r == nil || r.pool == nil {
		return nil, errors.New("job posting database is not configured")
	}
	rows, err := r.pool.Query(ctx, `SELECT `+postingColumns+
		` FROM job_postings WHERE company_id = $1::uuid AND deleted_at IS NULL ORDER BY created_at DESC, id`, companyID)
	if err != nil {
		return nil, fmt.Errorf("list company job postings: %w", err)
	}
	defer rows.Close()
	postings := make([]domain.Posting, 0)
	for rows.Next() {
		p, err := scanPosting(rows)
		if err != nil {
			return nil, fmt.Errorf("scan job posting: %w", err)
		}
		postings = append(postings, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read job postings: %w", err)
	}
	for i := range postings {
		if err := loadScreeningData(ctx, r.pool, &postings[i]); err != nil {
			return nil, err
		}
	}
	return postings, nil
}

func (r *PostingManagementRepository) Update(ctx context.Context, companyID, id string, p domain.Posting) (domain.Posting, error) {
	if r == nil || r.pool == nil {
		return domain.Posting{}, errors.New("job posting database is not configured")
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Posting{}, fmt.Errorf("begin posting update: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	row, err := scanPosting(tx.QueryRow(ctx, `UPDATE job_postings SET title=$3, description=$4, language=$5,
		currency=$6, seniority_level=$7, employment_type=$8, location=$9, salary_min=$10,
		salary_max=$11, salary_visible=$12, open_positions=$13, application_opening=$14,
		application_deadline=$15, max_applications=$16, updated_at=now()
		WHERE id=$1::uuid AND company_id=$2::uuid AND deleted_at IS NULL
		RETURNING `+postingColumns,
		id, companyID, strings.TrimSpace(p.Title), p.Description, p.Language, p.Currency,
		strings.TrimSpace(p.SeniorityLevel), p.EmploymentType, p.Location, p.SalaryMin,
		p.SalaryMax, p.SalaryVisible, p.OpenPositions, p.ApplicationOpening,
		p.ApplicationDeadline, p.MaxApplications))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Posting{}, domain.ErrPostingNotFound
	}
	if err != nil {
		return domain.Posting{}, fmt.Errorf("update job posting: %w", err)
	}
	if err := replaceScreeningData(ctx, tx, id, p); err != nil {
		return domain.Posting{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Posting{}, fmt.Errorf("commit posting update: %w", err)
	}
	return r.Get(ctx, companyID, id)
}

func (r *PostingManagementRepository) Publish(ctx context.Context, companyID, id string) (domain.Posting, error) {
	if r == nil || r.pool == nil {
		return domain.Posting{}, errors.New("job posting database is not configured")
	}
	row, err := scanPosting(r.pool.QueryRow(ctx, `UPDATE job_postings SET status='published',
		published_at=COALESCE(published_at, now()), archived_at=NULL, updated_at=now()
		WHERE id=$1::uuid AND company_id=$2::uuid AND deleted_at IS NULL AND status IN ('draft','archived')
		RETURNING `+postingColumns, id, companyID))
	if errors.Is(err, pgx.ErrNoRows) {
		if _, lookupErr := r.Get(ctx, companyID, id); errors.Is(lookupErr, domain.ErrPostingNotFound) {
			return domain.Posting{}, domain.ErrPostingNotFound
		}
		return domain.Posting{}, domain.ErrInvalidTransition
	}
	if err != nil {
		return domain.Posting{}, fmt.Errorf("publish job posting: %w", err)
	}
	if err := loadScreeningData(ctx, r.pool, &row); err != nil {
		return domain.Posting{}, err
	}
	return row, nil
}

func (r *PostingManagementRepository) Archive(ctx context.Context, companyID, id string) (domain.Posting, error) {
	if r == nil || r.pool == nil {
		return domain.Posting{}, errors.New("job posting database is not configured")
	}
	row, err := scanPosting(r.pool.QueryRow(ctx, `UPDATE job_postings SET status='archived', archived_at=now(), updated_at=now()
		WHERE id=$1::uuid AND company_id=$2::uuid AND deleted_at IS NULL AND status='published'
		RETURNING `+postingColumns, id, companyID))
	if errors.Is(err, pgx.ErrNoRows) {
		if _, lookupErr := r.Get(ctx, companyID, id); errors.Is(lookupErr, domain.ErrPostingNotFound) {
			return domain.Posting{}, domain.ErrPostingNotFound
		}
		return domain.Posting{}, domain.ErrInvalidTransition
	}
	if err != nil {
		return domain.Posting{}, fmt.Errorf("archive job posting: %w", err)
	}
	if err := loadScreeningData(ctx, r.pool, &row); err != nil {
		return domain.Posting{}, err
	}
	return row, nil
}

// SoftDelete marks the posting deleted. Admin deletion can target any company;
// manager ownership is checked in the HTTP/application boundary before calling.
func (r *PostingManagementRepository) SoftDelete(ctx context.Context, companyID, id string, systemAdmin bool) error {
	if r == nil || r.pool == nil {
		return errors.New("job posting database is not configured")
	}
	query := `UPDATE job_postings SET deleted_at=now(), updated_at=now() WHERE id=$1::uuid AND deleted_at IS NULL`
	args := []any{id}
	if !systemAdmin {
		query += ` AND company_id=$2::uuid`
		args = append(args, companyID)
	}
	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("soft delete job posting: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrPostingNotFound
	}
	return nil
}

func (r *PostingManagementRepository) CompanyCanPublish(ctx context.Context, companyID string) (bool, error) {
	if r == nil || r.pool == nil {
		return false, errors.New("job posting database is not configured")
	}
	var eligible bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM companies WHERE id=$1::uuid
		AND COALESCE(banned,FALSE)=FALSE AND deleted_at IS NULL AND deactivated_at IS NULL)`, companyID).Scan(&eligible)
	return eligible, err
}

func replaceScreeningData(ctx context.Context, tx pgx.Tx, postingID string, p domain.Posting) error {
	if _, err := tx.Exec(ctx, `DELETE FROM job_posting_requirements WHERE posting_id=$1::uuid`, postingID); err != nil {
		return fmt.Errorf("replace posting requirements: %w", err)
	}
	if _, err := tx.Exec(ctx, `DELETE FROM job_posting_amarat WHERE posting_id=$1::uuid`, postingID); err != nil {
		return fmt.Errorf("replace posting Amarat: %w", err)
	}
	for _, requirement := range p.Requirements {
		if _, err := tx.Exec(ctx, `INSERT INTO job_posting_requirements (id, posting_id, requirement_text, requirement_type)
			VALUES (gen_random_uuid(), $1::uuid, $2, $3)`, postingID, strings.TrimSpace(requirement.Text), requirement.Type); err != nil {
			return fmt.Errorf("insert posting requirement: %w", err)
		}
	}
	for _, amara := range p.Amarat {
		if _, err := tx.Exec(ctx, `INSERT INTO job_posting_amarat (id, posting_id, amara_title, amara_description)
			VALUES (gen_random_uuid(), $1::uuid, $2, $3)`, postingID, strings.TrimSpace(amara.Title), amara.Description); err != nil {
			return fmt.Errorf("insert posting Amara: %w", err)
		}
	}
	return nil
}

func loadScreeningData(ctx context.Context, db rowQuerier, p *domain.Posting) error {
	p.Requirements = make([]domain.Requirement, 0)
	rows, err := db.Query(ctx, `SELECT id::text, requirement_text, requirement_type FROM job_posting_requirements WHERE posting_id=$1::uuid ORDER BY id`, p.ID)
	if err != nil {
		return fmt.Errorf("load posting requirements: %w", err)
	}
	for rows.Next() {
		var requirement domain.Requirement
		if err := rows.Scan(&requirement.ID, &requirement.Text, &requirement.Type); err != nil {
			rows.Close()
			return fmt.Errorf("scan posting requirement: %w", err)
		}
		p.Requirements = append(p.Requirements, requirement)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("read posting requirements: %w", err)
	}
	rows.Close()
	p.Amarat = make([]domain.Amara, 0)
	rows, err = db.Query(ctx, `SELECT id::text, amara_title, amara_description FROM job_posting_amarat WHERE posting_id=$1::uuid ORDER BY id`, p.ID)
	if err != nil {
		return fmt.Errorf("load posting Amarat: %w", err)
	}
	for rows.Next() {
		var amara domain.Amara
		if err := rows.Scan(&amara.ID, &amara.Title, &amara.Description); err != nil {
			rows.Close()
			return fmt.Errorf("scan posting Amara: %w", err)
		}
		p.Amarat = append(p.Amarat, amara)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("read posting Amarat: %w", err)
	}
	rows.Close()
	return nil
}

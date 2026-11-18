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
	open_positions, application_opening, application_deadline, max_applications, finalist_target, status,
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
		&p.MaxApplications, &p.FinalistTarget, &status, &p.LastAutosavedAt, &p.PublishedAt, &p.ArchivedAt,
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
		 application_deadline, max_applications, finalist_target, status, created_at, updated_at)
		VALUES (gen_random_uuid(), $1::uuid, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, 'draft', now(), now())
		RETURNING `+postingColumns,
		companyID, strings.TrimSpace(p.Title), p.Description, p.Language, p.Currency,
		strings.TrimSpace(p.SeniorityLevel), p.EmploymentType, p.Location, p.SalaryMin,
		p.SalaryMax, p.SalaryVisible, p.OpenPositions, p.ApplicationOpening,
		p.ApplicationDeadline, p.MaxApplications, p.FinalistTarget))
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
	rows.Close()
	if err := loadScreeningDataBatch(ctx, r.pool, postings); err != nil {
		return nil, err
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
	_, err = scanPosting(tx.QueryRow(ctx, `UPDATE job_postings SET title=$3, description=$4, language=$5,
		currency=$6, seniority_level=$7, employment_type=$8, location=$9, salary_min=$10,
		salary_max=$11, salary_visible=$12, open_positions=$13, application_opening=$14,
		application_deadline=$15, max_applications=$16, finalist_target=$17, updated_at=now()
		WHERE id=$1::uuid AND company_id=$2::uuid AND deleted_at IS NULL
		RETURNING `+postingColumns,
		id, companyID, strings.TrimSpace(p.Title), p.Description, p.Language, p.Currency,
		strings.TrimSpace(p.SeniorityLevel), p.EmploymentType, p.Location, p.SalaryMin,
		p.SalaryMax, p.SalaryVisible, p.OpenPositions, p.ApplicationOpening,
		p.ApplicationDeadline, p.MaxApplications, p.FinalistTarget))
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
	owner := ""
	args := []any{id}
	if !systemAdmin {
		owner = ` AND company_id=$2::uuid`
		args = append(args, companyID)
	}
	// Deleting a posting also ends its hiring pipeline in the same statement:
	// active executions are cancelled so the worker stops admitting, assessing
	// and advancing candidates for a posting that no longer exists.
	var deleted int
	err := r.pool.QueryRow(ctx, `
		WITH deleted AS (
			UPDATE job_postings SET deleted_at=now(), updated_at=now()
			WHERE id=$1::uuid AND deleted_at IS NULL`+owner+` RETURNING id
		), cancelled AS (
			UPDATE pipeline_executions pe SET status='cancelled'
			FROM applications a JOIN deleted d ON d.id = a.posting_id
			WHERE pe.application_id = a.id AND pe.status = 'active'
			RETURNING pe.id
		)
		SELECT count(*) FROM deleted`, args...).Scan(&deleted)
	if err != nil {
		return fmt.Errorf("soft delete job posting: %w", err)
	}
	if deleted == 0 {
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
	postings := []domain.Posting{*p}
	if err := loadScreeningDataBatch(ctx, db, postings); err != nil {
		return err
	}
	*p = postings[0]
	return nil
}

// loadScreeningDataBatch fills requirements and Amarat for many postings with
// two queries in total, instead of two per posting.
func loadScreeningDataBatch(ctx context.Context, db rowQuerier, postings []domain.Posting) error {
	index := make(map[string]int, len(postings))
	ids := make([]string, len(postings))
	for i := range postings {
		postings[i].Requirements = make([]domain.Requirement, 0)
		postings[i].Amarat = make([]domain.Amara, 0)
		index[postings[i].ID] = i
		ids[i] = postings[i].ID
	}
	if len(ids) == 0 {
		return nil
	}
	rows, err := db.Query(ctx, `SELECT posting_id::text, id::text, requirement_text, requirement_type FROM job_posting_requirements WHERE posting_id = ANY($1::uuid[]) ORDER BY posting_id, id`, ids)
	if err != nil {
		return fmt.Errorf("load posting requirements: %w", err)
	}
	for rows.Next() {
		var posting string
		var requirement domain.Requirement
		if err := rows.Scan(&posting, &requirement.ID, &requirement.Text, &requirement.Type); err != nil {
			rows.Close()
			return fmt.Errorf("scan posting requirement: %w", err)
		}
		if i, ok := index[posting]; ok {
			postings[i].Requirements = append(postings[i].Requirements, requirement)
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("read posting requirements: %w", err)
	}
	rows, err = db.Query(ctx, `SELECT posting_id::text, id::text, amara_title, amara_description FROM job_posting_amarat WHERE posting_id = ANY($1::uuid[]) ORDER BY posting_id, id`, ids)
	if err != nil {
		return fmt.Errorf("load posting Amarat: %w", err)
	}
	for rows.Next() {
		var posting string
		var amara domain.Amara
		if err := rows.Scan(&posting, &amara.ID, &amara.Title, &amara.Description); err != nil {
			rows.Close()
			return fmt.Errorf("scan posting Amara: %w", err)
		}
		if i, ok := index[posting]; ok {
			postings[i].Amarat = append(postings[i].Amarat, amara)
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("read posting Amarat: %w", err)
	}
	return nil
}

package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CSCE-3701-Amara/amara/internal/cand/domain"
	"github.com/CSCE-3701-Amara/amara/internal/cand/ports"
)

// JobFeedRepository reads candidate-visible postings from PostgreSQL.
type JobFeedRepository struct {
	pool *pgxpool.Pool
}

func NewJobFeedRepository(pool *pgxpool.Pool) *JobFeedRepository {
	return &JobFeedRepository{pool: pool}
}

func (r *JobFeedRepository) ListPublished(ctx context.Context, filter ports.JobFeedQuery) (ports.JobFeedResult, error) {
	if r == nil || r.pool == nil {
		return ports.JobFeedResult{}, fmt.Errorf("candidate feed database is not configured")
	}
	args := []any{}
	where := []string{
		"status = 'published'", "deleted_at IS NULL", "published_at IS NOT NULL",
		"(application_opening IS NULL OR application_opening <= now())",
		"(application_deadline IS NULL OR application_deadline > now())",
	}
	add := func(expression string, value any) {
		args = append(args, value)
		where = append(where, fmt.Sprintf(expression, len(args)))
	}
	if filter.Title != "" {
		add("title ILIKE $%d", "%"+filter.Title+"%")
	}
	if filter.Language != "" {
		add("language = $%d", filter.Language)
	}
	if filter.Currency != "" {
		add("currency = $%d", filter.Currency)
	}
	if filter.SeniorityLevel != "" {
		add("seniority_level = $%d", filter.SeniorityLevel)
	}
	if filter.EmploymentType != "" {
		add("employment_type = $%d", filter.EmploymentType)
	}
	if filter.Location != "" {
		add("location ILIKE $%d", "%"+filter.Location+"%")
	}
	if filter.SalaryMin != nil {
		add("(salary_max IS NULL OR salary_max >= $%d)", *filter.SalaryMin)
	}
	if filter.SalaryMax != nil {
		add("(salary_min IS NULL OR salary_min <= $%d)", *filter.SalaryMax)
	}
	if filter.AfterPublishedAt != nil {
		args = append(args, *filter.AfterPublishedAt, filter.AfterPostingID)
		where = append(where, fmt.Sprintf("(published_at, id) < ($%d, $%d)", len(args)-1, len(args)))
	}
	args = append(args, filter.Limit+1)
	query := `SELECT id, company_id, title, description, language, currency, seniority_level, employment_type, location, salary_min, salary_max, salary_visible, open_positions, application_opening, application_deadline, max_applications, status, last_autosaved_at, published_at, archived_at, created_at, updated_at, deleted_at FROM job_postings WHERE ` + strings.Join(where, " AND ") + fmt.Sprintf(" ORDER BY published_at DESC, id DESC LIMIT $%d", len(args))
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return ports.JobFeedResult{}, fmt.Errorf("query candidate job feed: %w", err)
	}
	defer rows.Close()

	postings := make([]domain.FeedPosting, 0, filter.Limit+1)
	for rows.Next() {
		var posting domain.FeedPosting
		if err := rows.Scan(
			&posting.ID, &posting.CompanyID, &posting.Title, &posting.Description,
			&posting.Language, &posting.Currency, &posting.SeniorityLevel, &posting.EmploymentType,
			&posting.Location, &posting.SalaryMin, &posting.SalaryMax, &posting.SalaryVisible,
			&posting.OpenPositions, &posting.ApplicationOpening, &posting.ApplicationDeadline,
			&posting.MaxApplications, &posting.Status, &posting.LastAutosavedAt, &posting.PublishedAt,
			&posting.ArchivedAt, &posting.CreatedAt, &posting.UpdatedAt, &posting.DeletedAt,
		); err != nil {
			return ports.JobFeedResult{}, fmt.Errorf("scan candidate job feed row: %w", err)
		}
		postings = append(postings, posting)
	}
	if err := rows.Err(); err != nil {
		return ports.JobFeedResult{}, fmt.Errorf("read candidate job feed rows: %w", err)
	}
	result := ports.JobFeedResult{Postings: postings}
	if len(postings) > filter.Limit {
		result.HasMore = true
		result.Postings = postings[:filter.Limit]
	}
	return result, nil
}

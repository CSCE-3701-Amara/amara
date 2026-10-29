// Package application contains CAND use cases and orchestration.
package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FeedFilter contains the optional exact-match/range filters supported by the
// candidate job feed. Empty values are ignored.
type FeedFilter struct {
	Language       string
	Currency       string
	SeniorityLevel string
	EmploymentType string
	Location       string
	Title          string
	SalaryMin      *int32
	SalaryMax      *int32
	Limit          int
	Cursor         string
}

// FeedJob is a posting plus the link from which a candidate can begin applying.
// The application record and its status are deliberately not exposed here.
type FeedJob struct {
	ID                  pgtype.UUID      `json:"id"`
	CompanyID           pgtype.UUID      `json:"company_id"`
	Title               string           `json:"title"`
	Description         pgtype.Text      `json:"description"`
	Language            pgtype.Text      `json:"language"`
	Currency            pgtype.Text      `json:"currency"`
	SeniorityLevel      pgtype.Text      `json:"seniority_level"`
	EmploymentType      pgtype.Text      `json:"employment_type"`
	Location            pgtype.Text      `json:"location"`
	SalaryMin           pgtype.Int4      `json:"salary_min"`
	SalaryMax           pgtype.Int4      `json:"salary_max"`
	SalaryVisible       pgtype.Bool      `json:"salary_visible"`
	OpenPositions       pgtype.Int4      `json:"open_positions"`
	ApplicationOpening  *time.Time       `json:"application_opening"`
	ApplicationDeadline *time.Time       `json:"application_deadline"`
	MaxApplications     pgtype.Int4      `json:"max_applications"`
	Status              pgtype.Text      `json:"status"`
	LastAutosavedAt     pgtype.Timestamp `json:"last_autosaved_at"`
	PublishedAt         pgtype.Timestamp `json:"published_at"`
	ArchivedAt          pgtype.Timestamp `json:"archived_at"`
	CreatedAt           pgtype.Timestamp `json:"created_at"`
	UpdatedAt           pgtype.Timestamp `json:"updated_at"`
	DeletedAt           pgtype.Timestamp `json:"deleted_at"`
	ApplicationURL      string           `json:"application_url"`
}

// FeedPage is one page of postings and its optional next-page cursor.
type FeedPage struct {
	Items      []FeedJob `json:"items"`
	NextCursor string    `json:"next_cursor,omitempty"`
}

// JobFeed serves published, non-deleted postings from PostgreSQL.
type JobFeed struct{ pool *pgxpool.Pool }

func NewJobFeed(pool *pgxpool.Pool) *JobFeed { return &JobFeed{pool: pool} }

// List returns a cursor page ordered by PublishedAt DESC, then ID DESC. Cursor
// values encode the last row's timestamp and UUID as "RFC3339Nano|uuid".
func (s *JobFeed) List(ctx context.Context, filter FeedFilter) (FeedPage, error) {
	if s == nil || s.pool == nil {
		return FeedPage{}, fmt.Errorf("job feed database is not configured")
	}
	limit := filter.Limit
	if limit == 0 {
		limit = 20
	}
	if limit < 1 || limit > 100 {
		return FeedPage{}, fmt.Errorf("limit must be between 1 and 100")
	}
	// Expired postings transition to archived as part of feed reads, so status
	// and archived_at remain consistent even before a scheduled worker exists.
	if _, err := s.pool.Exec(ctx, `UPDATE job_postings SET status = 'archived', archived_at = COALESCE(archived_at, now()), updated_at = now() WHERE status = 'published' AND deleted_at IS NULL AND application_deadline IS NOT NULL AND application_deadline <= now()`); err != nil {
		return FeedPage{}, fmt.Errorf("archive expired job postings: %w", err)
	}

	args := []any{}
	where := []string{"status = 'published'", "deleted_at IS NULL", "published_at IS NOT NULL", "(application_opening IS NULL OR application_opening <= now())", "(application_deadline IS NULL OR application_deadline > now())"}
	add := func(expr string, value any) {
		args = append(args, value)
		where = append(where, fmt.Sprintf(expr, len(args)))
	}
	if v := strings.TrimSpace(filter.Title); v != "" {
		add("title ILIKE $%d", "%"+v+"%")
	}
	if v := strings.TrimSpace(filter.Language); v != "" {
		add("language = $%d", v)
	}
	if v := strings.TrimSpace(filter.Currency); v != "" {
		add("currency = $%d", v)
	}
	if v := strings.TrimSpace(filter.SeniorityLevel); v != "" {
		add("seniority_level = $%d", v)
	}
	if v := strings.TrimSpace(filter.EmploymentType); v != "" {
		add("employment_type = $%d", v)
	}
	if v := strings.TrimSpace(filter.Location); v != "" {
		add("location ILIKE $%d", "%"+v+"%")
	}
	if filter.SalaryMin != nil {
		add("(salary_max IS NULL OR salary_max >= $%d)", *filter.SalaryMin)
	}
	if filter.SalaryMax != nil {
		add("(salary_min IS NULL OR salary_min <= $%d)", *filter.SalaryMax)
	}
	if filter.Cursor != "" {
		parts := strings.SplitN(filter.Cursor, "|", 2)
		if len(parts) != 2 {
			return FeedPage{}, fmt.Errorf("invalid cursor")
		}
		timestamp, err := time.Parse(time.RFC3339Nano, parts[0])
		if err != nil {
			return FeedPage{}, fmt.Errorf("invalid cursor timestamp: %w", err)
		}
		var id pgtype.UUID
		if err := id.Scan(parts[1]); err != nil {
			return FeedPage{}, fmt.Errorf("invalid cursor UUID: %w", err)
		}
		args = append(args, timestamp, id)
		where = append(where, fmt.Sprintf("(published_at, id) < ($%d, $%d)", len(args)-1, len(args)))
	}
	args = append(args, limit+1)
	query := `SELECT id, company_id, title, description, language, currency, seniority_level, employment_type, location, salary_min, salary_max, salary_visible, open_positions, application_opening, application_deadline, max_applications, status, last_autosaved_at, published_at, archived_at, created_at, updated_at, deleted_at FROM job_postings WHERE ` + strings.Join(where, " AND ") + fmt.Sprintf(" ORDER BY published_at DESC, id DESC LIMIT $%d", len(args))
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return FeedPage{}, fmt.Errorf("query job feed: %w", err)
	}
	defer rows.Close()
	items := make([]FeedJob, 0, limit)
	for rows.Next() {
		var item FeedJob
		if err := rows.Scan(&item.ID, &item.CompanyID, &item.Title, &item.Description, &item.Language, &item.Currency, &item.SeniorityLevel, &item.EmploymentType, &item.Location, &item.SalaryMin, &item.SalaryMax, &item.SalaryVisible, &item.OpenPositions, &item.ApplicationOpening, &item.ApplicationDeadline, &item.MaxApplications, &item.Status, &item.LastAutosavedAt, &item.PublishedAt, &item.ArchivedAt, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt); err != nil {
			return FeedPage{}, fmt.Errorf("scan job feed row: %w", err)
		}
		item.ApplicationURL = "/v1/jobs/" + item.ID.String() + "/apply"
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return FeedPage{}, fmt.Errorf("read job feed rows: %w", err)
	}
	page := FeedPage{Items: items}
	if len(items) > limit {
		last := items[limit-1]
		page.Items = items[:limit]
		page.NextCursor = last.PublishedAt.Time.UTC().Format(time.RFC3339Nano) + "|" + last.ID.String()
	}
	return page, nil
}

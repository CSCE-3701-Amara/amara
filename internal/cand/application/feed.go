package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/cand/domain"
	"github.com/CSCE-3701-Amara/amara/internal/cand/ports"
	jobapplication "github.com/CSCE-3701-Amara/amara/internal/job/application"
)

// FeedFilter contains optional filters accepted by the candidate job feed.
type FeedFilter struct {
	Title          string
	Language       string
	Currency       string
	SeniorityLevel string
	EmploymentType string
	Location       string
	SalaryMin      *int32
	SalaryMax      *int32
	Limit          int
	Cursor         string
}

// FeedJob is a posting's public feed representation and its apply URL.
type FeedJob struct {
	ID                  string     `json:"id"`
	CompanyID           string     `json:"company_id"`
	Title               string     `json:"title"`
	Description         *string    `json:"description"`
	Language            *string    `json:"language"`
	Currency            *string    `json:"currency"`
	SeniorityLevel      *string    `json:"seniority_level"`
	EmploymentType      *string    `json:"employment_type"`
	Location            *string    `json:"location"`
	SalaryMin           *int32     `json:"salary_min"`
	SalaryMax           *int32     `json:"salary_max"`
	SalaryVisible       *bool      `json:"salary_visible"`
	OpenPositions       *int32     `json:"open_positions"`
	ApplicationOpening  *time.Time `json:"application_opening"`
	ApplicationDeadline *time.Time `json:"application_deadline"`
	MaxApplications     *int32     `json:"max_applications"`
	Status              string     `json:"status"`
	LastAutosavedAt     *time.Time `json:"last_autosaved_at"`
	PublishedAt         time.Time  `json:"published_at"`
	ArchivedAt          *time.Time `json:"archived_at"`
	CreatedAt           *time.Time `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at"`
	ApplicationURL      string     `json:"application_url"`
}

// FeedPage is one cursor-paginated page of postings.
type FeedPage struct {
	Items      []FeedJob `json:"items"`
	NextCursor string    `json:"next_cursor,omitempty"`
}

// JobFeed coordinates the CAND feed query with JOB's posting lifecycle
// contract. SQL and pool access stay inside the database adapters.
type JobFeed struct {
	repository ports.JobFeedRepository
	archiver   jobapplication.ExpiredPostingArchiver
}

func NewJobFeed(repository ports.JobFeedRepository, archiver jobapplication.ExpiredPostingArchiver) *JobFeed {
	return &JobFeed{repository: repository, archiver: archiver}
}

// List returns published postings ordered by PublishedAt DESC and ID DESC.
// Cursor values use RFC3339Nano and the posting UUID separated by a pipe.
func (s *JobFeed) List(ctx context.Context, filter FeedFilter) (FeedPage, error) {
	if s == nil || s.repository == nil || s.archiver == nil {
		return FeedPage{}, fmt.Errorf("job feed dependencies are not configured")
	}
	limit := filter.Limit
	if limit == 0 {
		limit = 20
	}
	if limit < 1 || limit > 100 {
		return FeedPage{}, fmt.Errorf("limit must be between 1 and 100")
	}

	query := ports.JobFeedQuery{
		Title: strings.TrimSpace(filter.Title), Language: strings.TrimSpace(filter.Language),
		Currency: strings.TrimSpace(filter.Currency), SeniorityLevel: strings.TrimSpace(filter.SeniorityLevel),
		EmploymentType: strings.TrimSpace(filter.EmploymentType), Location: strings.TrimSpace(filter.Location),
		SalaryMin: filter.SalaryMin, SalaryMax: filter.SalaryMax, Limit: limit,
	}
	if filter.Cursor != "" {
		parts := strings.SplitN(filter.Cursor, "|", 2)
		if len(parts) != 2 {
			return FeedPage{}, fmt.Errorf("invalid cursor")
		}
		publishedAt, err := time.Parse(time.RFC3339Nano, parts[0])
		if err != nil {
			return FeedPage{}, fmt.Errorf("invalid cursor timestamp: %w", err)
		}
		if !validUUID(parts[1]) {
			return FeedPage{}, fmt.Errorf("invalid cursor UUID")
		}
		query.AfterPublishedAt = &publishedAt
		query.AfterPostingID = parts[1]
	}

	if err := s.archiver.ArchiveExpired(ctx); err != nil {
		return FeedPage{}, fmt.Errorf("archive expired postings: %w", err)
	}
	result, err := s.repository.ListPublished(ctx, query)
	if err != nil {
		return FeedPage{}, fmt.Errorf("list published postings: %w", err)
	}

	page := FeedPage{Items: make([]FeedJob, 0, len(result.Postings))}
	for _, posting := range result.Postings {
		page.Items = append(page.Items, mapFeedPosting(posting))
	}
	if result.HasMore && len(page.Items) > 0 {
		last := page.Items[len(page.Items)-1]
		page.NextCursor = last.PublishedAt.UTC().Format(time.RFC3339Nano) + "|" + last.ID
	}
	return page, nil
}

func mapFeedPosting(posting domain.FeedPosting) FeedJob {
	return FeedJob{
		ID: posting.ID, CompanyID: posting.CompanyID, Title: posting.Title,
		Description: posting.Description, Language: posting.Language, Currency: posting.Currency,
		SeniorityLevel: posting.SeniorityLevel, EmploymentType: posting.EmploymentType,
		Location: posting.Location, SalaryMin: posting.SalaryMin, SalaryMax: posting.SalaryMax,
		SalaryVisible: posting.SalaryVisible, OpenPositions: posting.OpenPositions,
		ApplicationOpening: posting.ApplicationOpening, ApplicationDeadline: posting.ApplicationDeadline,
		MaxApplications: posting.MaxApplications, Status: posting.Status,
		LastAutosavedAt: posting.LastAutosavedAt, PublishedAt: posting.PublishedAt,
		ArchivedAt: posting.ArchivedAt, CreatedAt: posting.CreatedAt, UpdatedAt: posting.UpdatedAt,
		DeletedAt: posting.DeletedAt, ApplicationURL: "/v1/jobs/" + posting.ID + "/apply",
	}
}

func validUUID(value string) bool {
	if len(value) != 36 || value[8] != '-' || value[13] != '-' || value[18] != '-' || value[23] != '-' {
		return false
	}
	for i, r := range value {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue
		}
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return false
		}
	}
	return true
}

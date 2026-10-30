package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/job/domain"
	"github.com/CSCE-3701-Amara/amara/internal/job/ports"
)

// wraps every client-correctable validation failure

var ErrInvalidPostingListQuery = errors.New("invalid posting list query")

// PostingListFilter contains the company scope and optional filters accepted

type PostingListFilter struct {
	CompanyID string
	Title     string
	Status    string
	Limit     int
	Cursor    string
}

type CompanyPostingItem struct {
	ID                  string     `json:"id"`
	Title               string     `json:"title"`
	Status              string     `json:"status"`
	ApplicationCount    int64      `json:"application_count"`
	ApplicationDeadline *time.Time `json:"application_deadline"`
	PublishedAt         *time.Time `json:"published_at"`
	ArchivedAt          *time.Time `json:"archived_at"`
	CreatedAt           *time.Time `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
}

// CompanyPostingPage is one cursor-paginated page of company postings.
type CompanyPostingPage struct {
	Items      []CompanyPostingItem `json:"items"`
	NextCursor string               `json:"next_cursor,omitempty"`
}

// company posting with its state, application count, and deadline.
type CompanyPostingList struct {
	repository ports.CompanyPostingListRepository
	archiver   ExpiredPostingArchiver
}

// NewCompanyPostingList creates the company posting list service.
func NewCompanyPostingList(repository ports.CompanyPostingListRepository, archiver ExpiredPostingArchiver) *CompanyPostingList {
	return &CompanyPostingList{repository: repository, archiver: archiver}
}

func (s *CompanyPostingList) List(ctx context.Context, filter PostingListFilter) (CompanyPostingPage, error) {
	if s == nil || s.repository == nil || s.archiver == nil {
		return CompanyPostingPage{}, fmt.Errorf("company posting list dependencies are not configured")
	}
	if !validUUID(filter.CompanyID) {
		return CompanyPostingPage{}, fmt.Errorf("%w: company_id must be a UUID", ErrInvalidPostingListQuery)
	}
	status := strings.ToLower(strings.TrimSpace(filter.Status))
	if status != "" && !domain.ValidPostingStatus(status) {
		return CompanyPostingPage{}, fmt.Errorf("%w: status must be draft, published, or archived", ErrInvalidPostingListQuery)
	}
	limit := filter.Limit
	if limit == 0 {
		limit = 20
	}
	if limit < 1 || limit > 100 {
		return CompanyPostingPage{}, fmt.Errorf("%w: limit must be between 1 and 100", ErrInvalidPostingListQuery)
	}

	query := ports.CompanyPostingQuery{
		CompanyID: filter.CompanyID, Title: strings.TrimSpace(filter.Title), Status: status, Limit: limit,
	}
	if filter.Cursor != "" {
		parts := strings.SplitN(filter.Cursor, "|", 2)
		if len(parts) != 2 {
			return CompanyPostingPage{}, fmt.Errorf("%w: invalid cursor", ErrInvalidPostingListQuery)
		}
		createdAt, err := time.Parse(time.RFC3339Nano, parts[0])
		if err != nil {
			return CompanyPostingPage{}, fmt.Errorf("%w: invalid cursor timestamp", ErrInvalidPostingListQuery)
		}
		if !validUUID(parts[1]) {
			return CompanyPostingPage{}, fmt.Errorf("%w: invalid cursor UUID", ErrInvalidPostingListQuery)
		}
		// Columns are timestamp without time zone, so compare in UTC.
		createdAt = createdAt.UTC()
		query.AfterCreatedAt = &createdAt
		query.AfterPostingID = parts[1]
	}

	if err := s.archiver.ArchiveExpired(ctx); err != nil {
		return CompanyPostingPage{}, fmt.Errorf("archive expired postings: %w", err)
	}
	result, err := s.repository.ListByCompany(ctx, query)
	if err != nil {
		return CompanyPostingPage{}, fmt.Errorf("list company postings: %w", err)
	}

	page := CompanyPostingPage{Items: make([]CompanyPostingItem, 0, len(result.Postings))}
	for _, posting := range result.Postings {
		page.Items = append(page.Items, mapCompanyPosting(posting))
	}
	if result.HasMore && len(page.Items) > 0 {
		last := page.Items[len(page.Items)-1]
		page.NextCursor = CursorTime(last.CreatedAt).Format(time.RFC3339Nano) + "|" + last.ID
	}
	return page, nil
}

// CursorTime is the sort key for a posting's CreatedAt; NULL sorts as the zero
// time, matching the database adapter's COALESCE.
func CursorTime(createdAt *time.Time) time.Time {
	if createdAt == nil {
		return time.Time{}
	}
	return createdAt.UTC()
}

func mapCompanyPosting(posting domain.CompanyPosting) CompanyPostingItem {
	return CompanyPostingItem{
		ID: posting.ID, Title: posting.Title, Status: posting.Status,
		ApplicationCount: posting.ApplicationCount, ApplicationDeadline: posting.ApplicationDeadline,
		PublishedAt: posting.PublishedAt, ArchivedAt: posting.ArchivedAt,
		CreatedAt: posting.CreatedAt, UpdatedAt: posting.UpdatedAt,
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

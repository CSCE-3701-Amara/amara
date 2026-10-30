package application

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/cand/domain"
	"github.com/CSCE-3701-Amara/amara/internal/cand/ports"
)

type feedRepositoryStub struct {
	result ports.JobFeedResult
	err    error
	query  ports.JobFeedQuery
	calls  int
	events *[]string
}

func (s *feedRepositoryStub) ListPublished(_ context.Context, query ports.JobFeedQuery) (ports.JobFeedResult, error) {
	s.calls++
	s.query = query
	if s.events != nil {
		*s.events = append(*s.events, "list")
	}
	return s.result, s.err
}

type expiredArchiverStub struct {
	err    error
	calls  int
	events *[]string
}

func (s *expiredArchiverStub) ArchiveExpired(context.Context) error {
	s.calls++
	if s.events != nil {
		*s.events = append(*s.events, "archive")
	}
	return s.err
}

func TestJobFeedListAppliesFiltersAndMapsPosting(t *testing.T) {
	location := "Cairo"
	description := "Build reliable services"
	language := "en"
	currency := "EGP"
	seniority := "senior"
	employment := "full-time"
	minSalary, maxSalary, positions, maxApplications := int32(40000), int32(70000), int32(2), int32(50)
	salaryVisible := true
	publishedAt := time.Date(2026, time.June, 1, 12, 30, 0, 123456789, time.FixedZone("EET", 2*60*60))
	deadline := publishedAt.Add(30 * 24 * time.Hour)
	createdAt := publishedAt.Add(-24 * time.Hour)

	posting := domain.FeedPosting{
		ID: "4ce85a22-78c1-48bb-bbf6-1b2e06f6ad17", CompanyID: "1a1a77b2-32ed-45db-a437-7f66271e78af",
		Title: "Backend Engineer", Description: &description, Language: &language, Currency: &currency,
		SeniorityLevel: &seniority, EmploymentType: &employment, Location: &location,
		SalaryMin: &minSalary, SalaryMax: &maxSalary, SalaryVisible: &salaryVisible,
		OpenPositions: &positions, ApplicationDeadline: &deadline, MaxApplications: &maxApplications,
		Status: "published", PublishedAt: publishedAt, CreatedAt: &createdAt,
	}
	events := []string{}
	repository := &feedRepositoryStub{
		result: ports.JobFeedResult{Postings: []domain.FeedPosting{posting}, HasMore: true},
		events: &events,
	}
	archiver := &expiredArchiverStub{events: &events}
	service := NewJobFeed(repository, archiver)

	page, err := service.List(context.Background(), FeedFilter{
		Title: "  Backend ", Language: " en ", Currency: " EGP ", SeniorityLevel: " senior ",
		EmploymentType: " full-time ", Location: " Cairo ", SalaryMin: &minSalary, SalaryMax: &maxSalary,
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if archiver.calls != 1 || repository.calls != 1 {
		t.Fatalf("calls: archive=%d list=%d, want 1 each", archiver.calls, repository.calls)
	}
	if !reflect.DeepEqual(events, []string{"archive", "list"}) {
		t.Fatalf("operation order = %v, want [archive list]", events)
	}
	if repository.query.Title != "Backend" || repository.query.Language != "en" ||
		repository.query.Currency != "EGP" || repository.query.SeniorityLevel != "senior" ||
		repository.query.EmploymentType != "full-time" || repository.query.Location != "Cairo" {
		t.Fatalf("filters were not trimmed and forwarded: %+v", repository.query)
	}
	if repository.query.Limit != 20 {
		t.Fatalf("default limit = %d, want 20", repository.query.Limit)
	}
	if repository.query.SalaryMin == nil || *repository.query.SalaryMin != minSalary ||
		repository.query.SalaryMax == nil || *repository.query.SalaryMax != maxSalary {
		t.Fatalf("salary filters were not forwarded: %+v", repository.query)
	}
	if len(page.Items) != 1 {
		t.Fatalf("item count = %d, want 1", len(page.Items))
	}
	wantItem := FeedJob{
		ID: posting.ID, CompanyID: posting.CompanyID, Title: posting.Title, Description: posting.Description,
		Language: posting.Language, Currency: posting.Currency, SeniorityLevel: posting.SeniorityLevel,
		EmploymentType: posting.EmploymentType, Location: posting.Location, SalaryMin: posting.SalaryMin,
		SalaryMax: posting.SalaryMax, SalaryVisible: posting.SalaryVisible, OpenPositions: posting.OpenPositions,
		ApplicationDeadline: posting.ApplicationDeadline, MaxApplications: posting.MaxApplications,
		Status: posting.Status, PublishedAt: posting.PublishedAt, CreatedAt: posting.CreatedAt,
		ApplicationURL: "/v1/jobs/" + posting.ID + "/apply",
	}
	if !reflect.DeepEqual(page.Items[0], wantItem) {
		t.Errorf("mapped feed item = %#v, want %#v", page.Items[0], wantItem)
	}
	wantCursor := publishedAt.UTC().Format(time.RFC3339Nano) + "|" + posting.ID
	if page.NextCursor != wantCursor {
		t.Errorf("next cursor = %q, want %q", page.NextCursor, wantCursor)
	}
}

func TestJobFeedListDecodesCursorAndUsesRequestedLimit(t *testing.T) {
	cursorTime := time.Date(2026, time.July, 12, 9, 8, 7, 654321000, time.UTC)
	cursorID := "f95751c4-2a21-4e1a-8f32-c971ae9ac09d"
	repository := &feedRepositoryStub{}
	service := NewJobFeed(repository, &expiredArchiverStub{})

	_, err := service.List(context.Background(), FeedFilter{
		Limit: 7, Cursor: cursorTime.Format(time.RFC3339Nano) + "|" + cursorID,
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if repository.query.Limit != 7 {
		t.Errorf("limit = %d, want 7", repository.query.Limit)
	}
	if repository.query.AfterPublishedAt == nil || !repository.query.AfterPublishedAt.Equal(cursorTime) {
		t.Errorf("cursor timestamp = %v, want %v", repository.query.AfterPublishedAt, cursorTime)
	}
	if repository.query.AfterPostingID != cursorID {
		t.Errorf("cursor posting ID = %q, want %q", repository.query.AfterPostingID, cursorID)
	}
}

func TestJobFeedListRejectsInvalidLimitAndCursorBeforeDependencies(t *testing.T) {
	tests := []struct {
		name   string
		filter FeedFilter
	}{
		{name: "below minimum", filter: FeedFilter{Limit: -1}},
		{name: "above maximum", filter: FeedFilter{Limit: 101}},
		{name: "malformed cursor", filter: FeedFilter{Cursor: "not-a-cursor"}},
		{name: "invalid timestamp", filter: FeedFilter{Cursor: "not-a-time|4ce85a22-78c1-48bb-bbf6-1b2e06f6ad17"}},
		{name: "invalid UUID", filter: FeedFilter{Cursor: "2026-06-01T12:30:00Z|invalid"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &feedRepositoryStub{}
			archiver := &expiredArchiverStub{}
			service := NewJobFeed(repository, archiver)
			if _, err := service.List(context.Background(), test.filter); err == nil {
				t.Fatal("List() error = nil, want validation error")
			}
			if repository.calls != 0 || archiver.calls != 0 {
				t.Fatalf("invalid input reached dependencies: archive=%d list=%d", archiver.calls, repository.calls)
			}
		})
	}
}

func TestJobFeedListStopsWhenArchivingFails(t *testing.T) {
	archiveErr := errors.New("database unavailable")
	repository := &feedRepositoryStub{}
	archiver := &expiredArchiverStub{err: archiveErr}
	service := NewJobFeed(repository, archiver)

	_, err := service.List(context.Background(), FeedFilter{})
	if !errors.Is(err, archiveErr) {
		t.Fatalf("List() error = %v, want wrapped archive error", err)
	}
	if repository.calls != 0 {
		t.Fatalf("repository calls = %d, want 0 after archive failure", repository.calls)
	}
}

func TestJobFeedListPropagatesRepositoryError(t *testing.T) {
	repositoryErr := errors.New("query failed")
	repository := &feedRepositoryStub{err: repositoryErr}
	archiver := &expiredArchiverStub{}
	service := NewJobFeed(repository, archiver)

	_, err := service.List(context.Background(), FeedFilter{})
	if !errors.Is(err, repositoryErr) {
		t.Fatalf("List() error = %v, want wrapped repository error", err)
	}
	if archiver.calls != 1 || repository.calls != 1 {
		t.Fatalf("calls: archive=%d list=%d, want 1 each", archiver.calls, repository.calls)
	}
}

func TestJobFeedListRejectsMissingDependencies(t *testing.T) {
	var nilService *JobFeed
	if _, err := nilService.List(context.Background(), FeedFilter{}); err == nil {
		t.Fatal("nil service should return an error")
	}
	if _, err := NewJobFeed(nil, &expiredArchiverStub{}).List(context.Background(), FeedFilter{}); err == nil {
		t.Fatal("nil repository should return an error")
	}
	if _, err := NewJobFeed(&feedRepositoryStub{}, nil).List(context.Background(), FeedFilter{}); err == nil {
		t.Fatal("nil archiver should return an error")
	}
}

package ports

import (
	"context"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/cand/domain"
)

// JobFeedQuery contains validated filters and a decoded cursor for the
// candidate-facing feed query.
type JobFeedQuery struct {
	Title            string
	Language         string
	Currency         string
	SeniorityLevel   string
	EmploymentType   string
	Location         string
	SalaryMin        *int32
	SalaryMax        *int32
	Limit            int
	AfterPublishedAt *time.Time
	AfterPostingID   string
}

// JobFeedResult includes at most Limit postings and indicates whether one more
// row existed when the repository queried the database.
type JobFeedResult struct {
	Postings []domain.FeedPosting
	HasMore  bool
}

// JobFeedRepository is the CAND-owned read boundary for candidate job feeds.
type JobFeedRepository interface {
	ListPublished(ctx context.Context, query JobFeedQuery) (JobFeedResult, error)
}

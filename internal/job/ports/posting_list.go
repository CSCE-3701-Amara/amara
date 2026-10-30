package ports

import (
	"context"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/job/domain"
)

// contains validated filters and a decoded cursor for the company posting list query.

type CompanyPostingQuery struct {
	CompanyID      string
	Title          string
	Status         string
	Limit          int
	AfterCreatedAt *time.Time
	AfterPostingID string
}

type CompanyPostingResult struct {
	Postings []domain.CompanyPosting
	HasMore  bool
}

// CompanyPostingListRepository is the JOB-owned read boundary for a company's postings in every state.

type CompanyPostingListRepository interface {
	ListByCompany(ctx context.Context, query CompanyPostingQuery) (CompanyPostingResult, error)
}

package domain

import "time"

// FeedPosting is the candidate-facing read model for a published posting.
// Nullable database values are represented as pointers; application status is
// intentionally absent because applications are not part of the feed.
type FeedPosting struct {
	ID                  string
	CompanyID           string
	Title               string
	Description         *string
	Language            *string
	Currency            *string
	SeniorityLevel      *string
	EmploymentType      *string
	Location            *string
	SalaryMin           *int32
	SalaryMax           *int32
	SalaryVisible       *bool
	OpenPositions       *int32
	ApplicationOpening  *time.Time
	ApplicationDeadline *time.Time
	MaxApplications     *int32
	Status              string
	LastAutosavedAt     *time.Time
	PublishedAt         time.Time
	ArchivedAt          *time.Time
	CreatedAt           *time.Time
	UpdatedAt           *time.Time
	DeletedAt           *time.Time
}

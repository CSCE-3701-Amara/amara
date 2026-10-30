package domain

import "time"

// Posting states allowed.
const (
	PostingStatusDraft     = "draft"
	PostingStatusPublished = "published"
	PostingStatusArchived  = "archived"
)

// ValidPostingStatus reports whether status is one of the JOB posting states.
func ValidPostingStatus(status string) bool {
	switch status {
	case PostingStatusDraft, PostingStatusPublished, PostingStatusArchived:
		return true
	}
	return false
}

type CompanyPosting struct {
	ID                  string
	Title               string
	Status              string
	ApplicationCount    int64
	ApplicationDeadline *time.Time
	PublishedAt         *time.Time
	ArchivedAt          *time.Time
	CreatedAt           *time.Time
	UpdatedAt           *time.Time
}

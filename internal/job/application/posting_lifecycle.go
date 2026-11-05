package application

import (
	"context"
	"fmt"

	"github.com/CSCE-3701-Amara/amara/internal/job/ports"
)

// ExpiredPostingArchiver is JOB's application contract for archiving postings
// whose application deadline has passed.
type ExpiredPostingArchiver interface {
	ArchiveExpired(ctx context.Context) error
}

// PostingLifecycle applies JOB-owned posting lifecycle transitions.
type PostingLifecycle struct {
	repository ports.PostingLifecycleRepository
}

func NewPostingLifecycle(repository ports.PostingLifecycleRepository) *PostingLifecycle {
	return &PostingLifecycle{repository: repository}
}

// ArchiveExpired transitions published postings past their deadline to
// archived. The database adapter performs this update in the JOB boundary.
func (s *PostingLifecycle) ArchiveExpired(ctx context.Context) error {
	if s == nil || s.repository == nil {
		return fmt.Errorf("job posting lifecycle repository is not configured")
	}
	if err := s.repository.ArchiveExpired(ctx); err != nil {
		return fmt.Errorf("archive expired job postings: %w", err)
	}
	return nil
}

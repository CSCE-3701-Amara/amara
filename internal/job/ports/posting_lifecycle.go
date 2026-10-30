package ports

import "context"

// PostingLifecycleRepository owns persistence for job posting state changes.
type PostingLifecycleRepository interface {
	ArchiveExpired(ctx context.Context) error
}

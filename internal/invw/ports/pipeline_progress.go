package ports

import (
	"context"

	"github.com/CSCE-3701-Amara/amara/internal/invw/domain"
)

// CandidatePipelineProgressRepository returns only a published pipeline
// assigned to the specified candidate-owned application.
type CandidatePipelineProgressRepository interface {
	GetByCandidateApplication(ctx context.Context, candidateID, applicationID string) (domain.CandidatePipelineProgress, error)
}

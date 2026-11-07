package ports

import (
	"context"

	"github.com/CSCE-3701-Amara/amara/internal/cand/domain"
)

// AssessmentResultRepository reads submitted assessments and only returns
// outcomes or feedback after their respective disclosure timestamps are set.
type AssessmentResultRepository interface {
	ListByCandidateApplication(ctx context.Context, candidateID, applicationID string) ([]domain.AssessmentResult, bool, error)
}

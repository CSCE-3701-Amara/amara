package ports

import (
	"context"
	"encoding/json"

	"github.com/CSCE-3701-Amara/amara/internal/cand/domain"
)

// AssessmentProgressRepository scopes every operation to the authenticated
// candidate and the candidate's application.
type AssessmentProgressRepository interface {
	ListByCandidateApplication(ctx context.Context, candidateID, applicationID string) ([]domain.AssessmentProgress, bool, error)
	SaveProgress(ctx context.Context, candidateID, applicationID, assessmentID string, progress json.RawMessage) (domain.AssessmentProgress, error)
}

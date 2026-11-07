package application

import (
	"context"
	"fmt"
	"github.com/CSCE-3701-Amara/amara/internal/job/domain"
	"github.com/CSCE-3701-Amara/amara/internal/job/ports"
)

// ApplicationAdmission accepts a candidate submission and runs initial screening.
type ApplicationAdmission struct{ repository ports.ApplicationRepository }

func NewApplicationAdmission(repository ports.ApplicationRepository) *ApplicationAdmission {
	return &ApplicationAdmission{repository: repository}
}
func (s *ApplicationAdmission) Apply(ctx context.Context, postingID, candidateID string) (domain.Application, error) {
	if s == nil || s.repository == nil {
		return domain.Application{}, fmt.Errorf("application admission is not configured")
	}
	return s.repository.Apply(ctx, postingID, candidateID)
}

package ports

import (
	"context"
	"github.com/CSCE-3701-Amara/amara/internal/job/domain"
)

// ApplicationRepository owns admission and initial screening persistence.
type ApplicationRepository interface {
	Apply(context.Context, string, string) (domain.Application, error)
}

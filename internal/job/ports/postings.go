package ports

import (
	"context"

	"github.com/CSCE-3701-Amara/amara/internal/job/domain"
)

type PostingManagementRepository interface {
	Create(context.Context, string, domain.Posting) (domain.Posting, error)
	Get(context.Context, string, string) (domain.Posting, error)
	ListByCompany(context.Context, string) ([]domain.Posting, error)
	Update(context.Context, string, string, domain.Posting) (domain.Posting, error)
	Publish(context.Context, string, string) (domain.Posting, error)
	Archive(context.Context, string, string) (domain.Posting, error)
	SoftDelete(context.Context, string, string, bool) error
	CompanyCanPublish(context.Context, string) (bool, error)
}

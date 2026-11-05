package application

import (
	"context"
	"fmt"

	"github.com/CSCE-3701-Amara/amara/internal/job/domain"
	"github.com/CSCE-3701-Amara/amara/internal/job/ports"
)

type PostingManagement struct {
	repository ports.PostingManagementRepository
}

func NewPostingManagement(repository ports.PostingManagementRepository) *PostingManagement {
	return &PostingManagement{repository: repository}
}

func (s *PostingManagement) Create(ctx context.Context, companyID string, posting domain.Posting) (domain.Posting, error) {
	posting = posting.Normalized()
	if err := posting.ValidateForSave(); err != nil {
		return domain.Posting{}, err
	}
	return s.repository.Create(ctx, companyID, posting)
}

func (s *PostingManagement) Get(ctx context.Context, companyID, id string) (domain.Posting, error) {
	return s.repository.Get(ctx, companyID, id)
}

func (s *PostingManagement) List(ctx context.Context, companyID string) ([]domain.Posting, error) {
	return s.repository.ListByCompany(ctx, companyID)
}

func (s *PostingManagement) Update(ctx context.Context, companyID, id string, posting domain.Posting) (domain.Posting, error) {
	posting = posting.Normalized()
	if err := posting.ValidateForSave(); err != nil {
		return domain.Posting{}, err
	}
	return s.repository.Update(ctx, companyID, id, posting)
}

func (s *PostingManagement) Publish(ctx context.Context, companyID, id string) (domain.Posting, error) {
	posting, err := s.repository.Get(ctx, companyID, id)
	if err != nil {
		return domain.Posting{}, err
	}
	if err := posting.ValidateForSave(); err != nil {
		return domain.Posting{}, err
	}
	canPublish, err := s.repository.CompanyCanPublish(ctx, companyID)
	if err != nil {
		return domain.Posting{}, fmt.Errorf("check company publishing eligibility: %w", err)
	}
	if !canPublish {
		return domain.Posting{}, fmt.Errorf("%w: company is inactive or banned", domain.ErrInvalidTransition)
	}
	return s.repository.Publish(ctx, companyID, id)
}

func (s *PostingManagement) Archive(ctx context.Context, companyID, id string) (domain.Posting, error) {
	return s.repository.Archive(ctx, companyID, id)
}

func (s *PostingManagement) Delete(ctx context.Context, companyID, id string, systemAdmin bool) error {
	return s.repository.SoftDelete(ctx, companyID, id, systemAdmin)
}

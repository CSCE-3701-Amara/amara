package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type AccountActionService struct {
	repository ports.AccountActionRepository
}

func NewAccountActionService(repository ports.AccountActionRepository) (*AccountActionService, error) {
	if repository == nil {
		return nil, errors.New("account action repository must not be nil")
	}
	return &AccountActionService{repository: repository}, nil
}

func (s *AccountActionService) Perform(ctx context.Context, accountID, sessionToken string, accountType domain.AccountType, action domain.AccountAction) error {
	if err := action.Validate(); err != nil {
		return err
	}
	if accountID == "" || sessionToken == "" || !validSelfServiceAccountType(accountType) {
		return domain.ErrInvalidSession
	}
	if err := s.repository.PerformAccountAction(ctx, accountID, sessionToken, accountType, action); err != nil {
		return fmt.Errorf("perform account action: %w", err)
	}
	return nil
}

func validSelfServiceAccountType(accountType domain.AccountType) bool {
	return accountType == domain.AccountTypeCandidate ||
		accountType == domain.AccountTypeInterviewer ||
		accountType == domain.AccountTypeCompany
}

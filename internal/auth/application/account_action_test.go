package application

import (
	"context"
	"errors"
	"testing"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

type accountActionRepository struct {
	accountID   string
	token       string
	accountType domain.AccountType
	action      domain.AccountAction
	err         error
}

func (r *accountActionRepository) PerformAccountAction(_ context.Context, accountID, token string, accountType domain.AccountType, action domain.AccountAction) error {
	r.accountID, r.token, r.accountType, r.action = accountID, token, accountType, action
	return r.err
}

func TestAccountActionServicePerformsAction(t *testing.T) {
	repository := &accountActionRepository{}
	service, err := NewAccountActionService(repository)
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Perform(context.Background(), "candidate-id", "session-token", domain.AccountTypeCandidate, domain.AccountActionDelete); err != nil {
		t.Fatal(err)
	}
	if repository.accountID != "candidate-id" || repository.token != "session-token" ||
		repository.accountType != domain.AccountTypeCandidate || repository.action != domain.AccountActionDelete {
		t.Fatalf("unexpected action: %+v", repository)
	}
}

func TestAccountActionServiceValidatesInputAndMapsFailure(t *testing.T) {
	repository := &accountActionRepository{err: errors.New("storage down")}
	service, err := NewAccountActionService(repository)
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Perform(context.Background(), "candidate-id", "session-token", domain.AccountTypeCandidate, "invalid"); !errors.Is(err, domain.ErrInvalidAccountAction) {
		t.Fatalf("invalid action error = %v", err)
	}
	if err := service.Perform(context.Background(), "candidate-id", "session-token", domain.AccountTypeCandidate, domain.AccountActionDeactivate); err == nil {
		t.Fatal("expected persistence error")
	}
	if err := service.Perform(context.Background(), "account-id", "session-token", "admin", domain.AccountActionDeactivate); !errors.Is(err, domain.ErrInvalidSession) {
		t.Fatalf("unsupported account type error = %v", err)
	}
}

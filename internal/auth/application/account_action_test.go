package application

import (
	"context"
	"errors"
	"testing"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

type accountActionRepository struct {
	accountID string
	token     string
	action    domain.AccountAction
	err       error
}

func (r *accountActionRepository) PerformAccountAction(_ context.Context, accountID, token string, action domain.AccountAction) error {
	r.accountID, r.token, r.action = accountID, token, action
	return r.err
}

func TestAccountActionServicePerformsAction(t *testing.T) {
	repository := &accountActionRepository{}
	service, err := NewAccountActionService(repository)
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Perform(context.Background(), "candidate-id", "session-token", domain.AccountActionDelete); err != nil {
		t.Fatal(err)
	}
	if repository.accountID != "candidate-id" || repository.token != "session-token" ||
		repository.action != domain.AccountActionDelete {
		t.Fatalf("unexpected action: %+v", repository)
	}
}

func TestAccountActionServiceValidatesInputAndMapsFailure(t *testing.T) {
	repository := &accountActionRepository{err: errors.New("storage down")}
	service, err := NewAccountActionService(repository)
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Perform(context.Background(), "candidate-id", "session-token", "invalid"); !errors.Is(err, domain.ErrInvalidAccountAction) {
		t.Fatalf("invalid action error = %v", err)
	}
	if err := service.Perform(context.Background(), "candidate-id", "session-token", domain.AccountActionDeactivate); err == nil {
		t.Fatal("expected persistence error")
	}
}

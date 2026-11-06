package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

func validCompanyRegistration() domain.CompanyRegistration {
	return domain.CompanyRegistration{
		CompanyName: "Amara Ltd", RegistrationNumber: "REG-123",
		Email: "OWNER@example.com", Password: "safe password",
	}
}

func newCompanyService(t *testing.T, repository *fakeRepository, checker fakeChecker) *CompanyRegistrationService {
	t.Helper()
	service, err := NewCompanyRegistrationService(repository, checker)
	if err != nil {
		t.Fatal(err)
	}
	service.now = func() time.Time { return time.Date(2026, 10, 30, 20, 0, 0, 0, time.UTC) }
	return service
}

func TestNewCompanyRegistrationServiceRequiresDependencies(t *testing.T) {
	if _, err := NewCompanyRegistrationService(nil, fakeChecker{}); err == nil {
		t.Fatal("expected repository error")
	}
	if _, err := NewCompanyRegistrationService(&fakeRepository{}, nil); err == nil {
		t.Fatal("expected checker error")
	}
}

func TestCompanyRegistrationServiceRegister(t *testing.T) {
	repository := &fakeRepository{account: domain.Account{
		ID: "request-id", Email: "owner@example.com", Status: domain.PendingReview,
	}}
	account, err := newCompanyService(t, repository, fakeChecker{}).Register(context.Background(), validCompanyRegistration())
	if err != nil {
		t.Fatal(err)
	}
	if account.ID != "request-id" || account.Status != domain.PendingReview {
		t.Fatalf("unexpected account: %+v", account)
	}
}

func TestCompanyRegistrationServiceFailures(t *testing.T) {
	tests := []struct {
		name       string
		repository *fakeRepository
		checker    fakeChecker
		want       error
	}{
		{"invalid input", &fakeRepository{}, fakeChecker{}, domain.ErrInvalidRegistration},
		{"email lookup", &fakeRepository{emailErr: errors.New("db down")}, fakeChecker{}, nil},
		{"duplicate email", &fakeRepository{emailExists: true}, fakeChecker{}, domain.ErrEmailAlreadyExists},
		{"number lookup", &fakeRepository{registrationNumberErr: errors.New("db down")}, fakeChecker{}, nil},
		{"duplicate number", &fakeRepository{registrationNumberExists: true}, fakeChecker{}, domain.ErrCompanyRegistrationExists},
		{"password lookup", &fakeRepository{}, fakeChecker{err: errors.New("service down")}, nil},
		{"compromised password", &fakeRepository{}, fakeChecker{compromised: true}, domain.ErrInvalidRegistration},
		{"create request", &fakeRepository{companyRegistrationErr: errors.New("insert failed")}, fakeChecker{}, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration := validCompanyRegistration()
			if test.name == "invalid input" {
				registration = domain.CompanyRegistration{}
			}
			_, err := newCompanyService(t, test.repository, test.checker).Register(context.Background(), registration)
			if test.want != nil && !errors.Is(err, test.want) {
				t.Fatalf("expected %v, got %v", test.want, err)
			}
			if test.want == nil && err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

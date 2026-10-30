package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type fakeRepository struct {
	emailExists              bool
	emailErr                 error
	companyID                string
	companyErr               error
	account                  domain.Account
	createErr                error
	createdCandidate         bool
	createdInterviewer       bool
	challenge                ports.VerificationChallenge
	challengeErr             error
	registrationNumberExists bool
	registrationNumberErr    error
	companyRegistrationErr   error
}

func (f *fakeRepository) EmailExists(context.Context, string) (bool, error) {
	return f.emailExists, f.emailErr
}
func (f *fakeRepository) FindApprovedCompanyBySecret(context.Context, string) (string, error) {
	return f.companyID, f.companyErr
}
func (f *fakeRepository) CreateCandidate(context.Context, string, string, []byte, []byte, time.Time) (domain.Account, error) {
	f.createdCandidate = true
	return f.account, f.createErr
}
func (f *fakeRepository) CreateInterviewer(context.Context, string, string, []byte, []byte, string, time.Time) (domain.Account, error) {
	f.createdInterviewer = true
	return f.account, f.createErr
}
func (f *fakeRepository) SaveVerificationChallenge(_ context.Context, challenge ports.VerificationChallenge) error {
	f.challenge = challenge
	return f.challengeErr
}
func (f *fakeRepository) RegistrationNumberExists(context.Context, string) (bool, error) {
	return f.registrationNumberExists, f.registrationNumberErr
}
func (f *fakeRepository) CreateCompanyRegistration(context.Context, ports.CompanyRegistrationRecord) (domain.Account, error) {
	return f.account, f.companyRegistrationErr
}

type fakeChecker struct {
	compromised bool
	err         error
}

func (f fakeChecker) IsCompromised(context.Context, string) (bool, error) {
	return f.compromised, f.err
}

type fakeSender struct {
	email string
	code  string
	err   error
}

func (f *fakeSender) SendVerificationCode(_ context.Context, email, code string) error {
	f.email, f.code = email, code
	return f.err
}

func validRegistration() domain.Registration {
	return domain.Registration{
		AccountType: domain.AccountTypeCandidate, Email: "USER@example.com",
		Password: "safe password", FirstName: "Ada", LastName: "Lovelace",
	}
}

func newService(t *testing.T, repository *fakeRepository, checker fakeChecker, sender *fakeSender) *RegistrationService {
	t.Helper()
	service, err := NewRegistrationService(repository, checker, sender)
	if err != nil {
		t.Fatal(err)
	}
	service.now = func() time.Time { return time.Date(2026, 10, 30, 20, 0, 0, 0, time.UTC) }
	return service
}

func TestNewRegistrationServiceRequiresDependencies(t *testing.T) {
	sender := &fakeSender{}
	if _, err := NewRegistrationService(nil, fakeChecker{}, sender); err == nil {
		t.Fatal("expected nil repository error")
	}
	if _, err := NewRegistrationService(&fakeRepository{}, nil, sender); err == nil {
		t.Fatal("expected nil checker error")
	}
	if _, err := NewRegistrationService(&fakeRepository{}, fakeChecker{}, nil); err == nil {
		t.Fatal("expected nil sender error")
	}
}

func TestRegisterCandidate(t *testing.T) {
	repository := &fakeRepository{account: domain.Account{ID: "candidate-id", Type: domain.AccountTypeCandidate, Email: "user@example.com", Status: domain.PendingVerification, CreatedAt: "2026-10-30T20:00:00Z"}}
	sender := &fakeSender{}
	account, err := newService(t, repository, fakeChecker{}, sender).Register(context.Background(), validRegistration())
	if err != nil {
		t.Fatal(err)
	}
	if account.ID != "candidate-id" || !repository.createdCandidate || repository.createdInterviewer {
		t.Fatalf("unexpected candidate registration: %+v", account)
	}
	if sender.email != "user@example.com" || len(sender.code) != 64 {
		t.Fatalf("unexpected verification message: email=%q code length=%d", sender.email, len(sender.code))
	}
	if repository.challenge.AccountID != account.ID || len(repository.challenge.CodeHash) != 32 {
		t.Fatalf("unexpected challenge: %+v", repository.challenge)
	}
}

func TestRegisterInterviewer(t *testing.T) {
	repository := &fakeRepository{
		companyID: "123e4567-e89b-12d3-a456-426614174000",
		account:   domain.Account{ID: "interviewer-id", Type: domain.AccountTypeInterviewer, Email: "interviewer@example.com", Status: domain.PendingVerification, CreatedAt: "2026-10-30T20:00:00Z"},
	}
	registration := validRegistration()
	registration.AccountType = domain.AccountTypeInterviewer
	registration.Email = "interviewer@example.com"
	registration.CompanySecret = "secret"
	if _, err := newService(t, repository, fakeChecker{}, &fakeSender{}).Register(context.Background(), registration); err != nil {
		t.Fatal(err)
	}
	if !repository.createdInterviewer {
		t.Fatal("expected interviewer account creation")
	}
}

func TestRegisterFailures(t *testing.T) {
	t.Run("invalid registration", func(t *testing.T) {
		_, err := newService(t, &fakeRepository{}, fakeChecker{}, &fakeSender{}).Register(context.Background(), domain.Registration{})
		if !errors.Is(err, domain.ErrInvalidRegistration) {
			t.Fatalf("expected invalid registration, got %v", err)
		}
	})
	t.Run("email exists", func(t *testing.T) {
		_, err := newService(t, &fakeRepository{emailExists: true}, fakeChecker{}, &fakeSender{}).Register(context.Background(), validRegistration())
		if !errors.Is(err, domain.ErrEmailAlreadyExists) {
			t.Fatalf("expected duplicate email, got %v", err)
		}
	})
	t.Run("email lookup error", func(t *testing.T) {
		_, err := newService(t, &fakeRepository{emailErr: errors.New("db down")}, fakeChecker{}, &fakeSender{}).Register(context.Background(), validRegistration())
		if err == nil {
			t.Fatal("expected lookup error")
		}
	})
	t.Run("compromise checker error", func(t *testing.T) {
		_, err := newService(t, &fakeRepository{}, fakeChecker{err: errors.New("service down")}, &fakeSender{}).Register(context.Background(), validRegistration())
		if err == nil {
			t.Fatal("expected checker error")
		}
	})
	t.Run("compromised password", func(t *testing.T) {
		_, err := newService(t, &fakeRepository{}, fakeChecker{compromised: true}, &fakeSender{}).Register(context.Background(), validRegistration())
		if !errors.Is(err, domain.ErrInvalidRegistration) {
			t.Fatalf("expected invalid registration, got %v", err)
		}
	})
	t.Run("company not found", func(t *testing.T) {
		registration := validRegistration()
		registration.AccountType = domain.AccountTypeInterviewer
		registration.CompanySecret = "secret"
		_, err := newService(t, &fakeRepository{companyErr: domain.ErrCompanyNotFound}, fakeChecker{}, &fakeSender{}).Register(context.Background(), registration)
		if !errors.Is(err, domain.ErrCompanyNotFound) {
			t.Fatalf("expected company error, got %v", err)
		}
	})
	t.Run("create error", func(t *testing.T) {
		_, err := newService(t, &fakeRepository{createErr: errors.New("insert failed"), account: domain.Account{ID: "id"}}, fakeChecker{}, &fakeSender{}).Register(context.Background(), validRegistration())
		if err == nil {
			t.Fatal("expected create error")
		}
	})
	t.Run("challenge error", func(t *testing.T) {
		repository := &fakeRepository{account: domain.Account{ID: "id"}, challengeErr: errors.New("insert failed")}
		_, err := newService(t, repository, fakeChecker{}, &fakeSender{}).Register(context.Background(), validRegistration())
		if err == nil {
			t.Fatal("expected challenge error")
		}
	})
	t.Run("sender error", func(t *testing.T) {
		repository := &fakeRepository{account: domain.Account{ID: "id"}}
		_, err := newService(t, repository, fakeChecker{}, &fakeSender{err: errors.New("mail failed")}).Register(context.Background(), validRegistration())
		if err == nil {
			t.Fatal("expected sender error")
		}
	})
}

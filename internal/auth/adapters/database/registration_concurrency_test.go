package database

import (
	"context"
	"errors"
	authapplication "github.com/CSCE-3701-Amara/amara/internal/auth/application"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
	"github.com/CSCE-3701-Amara/amara/internal/testutil"
)

func TestConcurrentRegistrationReservesEmailAcrossAccountTypesPostgres(t *testing.T) {
	pool := testutil.Postgres(t)
	ctx := t.Context()
	repo, err := NewRegistrationRepository(pool)
	if err != nil {
		t.Fatal(err)
	}
	const company = "00000000-0000-0000-0000-000000000410"
	if _, err := pool.Exec(ctx, `INSERT INTO companies(id,email,password_hash,salt,legal_name,approved)
 VALUES($1::uuid,'association@example.test','','','Association',true)`, company); err != nil {
		t.Fatal(err)
	}
	const email = "shared@example.test"
	create := []func() error{
		func() error {
			_, err := repo.CreateCandidate(ctx, "Candidate", email, []byte{}, []byte{}, time.Now())
			return err
		},
		func() error {
			_, err := repo.CreateInterviewer(ctx, "Interviewer", strings.ToUpper(email), []byte{}, []byte{}, company, time.Now())
			return err
		},
		func() error {
			_, err := repo.CreateCompanyRegistration(ctx, ports.CompanyRegistrationRecord{CompanyName: "Company", Email: email, RegistrationNumber: "12345", PasswordHash: []byte{}, Salt: []byte{}, CreatedAt: time.Now()})
			return err
		},
	}
	// All requests have passed the application's preflight before any insert.
	for range create {
		if exists, err := repo.EmailExists(ctx, email); err != nil || exists {
			t.Fatalf("preflight: exists=%v, error=%v", exists, err)
		}
	}
	var wg sync.WaitGroup
	errs := make(chan error, len(create))
	for _, fn := range create {
		wg.Go(func() { errs <- fn() })
	}
	wg.Wait()
	close(errs)
	successes, duplicates := 0, 0
	for err := range errs {
		if err == nil {
			successes++
		} else if errors.Is(err, domain.ErrEmailAlreadyExists) {
			duplicates++
		} else {
			t.Errorf("registration error: %v", err)
		}
	}
	if successes != 1 || duplicates != 2 {
		t.Fatalf("created %d accounts, rejected %d duplicates; want one email owner", successes, duplicates)
	}
}

func TestDuplicateCompanyNumberReturnsNumberConflictPostgres(t *testing.T) {
	pool := testutil.Postgres(t)
	repo, err := NewRegistrationRepository(pool)
	if err != nil {
		t.Fatal(err)
	}
	record := ports.CompanyRegistrationRecord{CompanyName: "Company", RegistrationNumber: "REG-12345", Email: "first@example.test", PasswordHash: []byte{}, Salt: []byte{}, CreatedAt: time.Now()}
	if _, err := repo.CreateCompanyRegistration(t.Context(), record); err != nil {
		t.Fatal(err)
	}
	record.Email = "second@example.test"
	if _, err := repo.CreateCompanyRegistration(t.Context(), record); !errors.Is(err, domain.ErrCompanyRegistrationExists) {
		t.Fatalf("duplicate company number error=%v, want company registration conflict", err)
	}
}

type registrationDelivery struct{ code string }

func (*registrationDelivery) IsCompromised(context.Context, string) (bool, error) { return false, nil }
func (d *registrationDelivery) SendVerificationCode(_ context.Context, _ string, code string) error {
	d.code = code
	return nil
}

func TestPendingRegistrationCredentialRoundTripPostgres(t *testing.T) {
	pool := testutil.Postgres(t)
	ctx := t.Context()
	repo, err := NewRegistrationRepository(pool)
	if err != nil {
		t.Fatal(err)
	}
	delivery := &registrationDelivery{}
	service, err := authapplication.NewRegistrationService(repo, delivery, delivery)
	if err != nil {
		t.Fatal(err)
	}
	registration := domain.Registration{AccountType: domain.AccountTypeCandidate, Email: "pending@example.test", Password: "original password", FirstName: "Pending", LastName: "Candidate"}
	original, err := service.Register(ctx, registration)
	if err != nil {
		t.Fatal(err)
	}
	originalCode := delivery.code
	registration.Password = "another password"
	if _, err := service.Register(ctx, registration); !errors.Is(err, domain.ErrEmailAlreadyExists) {
		t.Fatalf("wrong-password retry: %v", err)
	}
	if delivery.code != originalCode {
		t.Fatal("unauthorized retry replaced the challenge")
	}
	registration.Password = "original password"
	retry, err := service.Register(ctx, registration)
	if err != nil || retry.ID != original.ID || delivery.code == originalCode {
		t.Fatalf("valid resend failed: %+v, %v", retry, err)
	}
	verification, err := authapplication.NewVerificationService(repo)
	if err != nil {
		t.Fatal(err)
	}
	account, err := verification.Verify(ctx, registration.Email, delivery.code)
	if err != nil || account.Status != "active" {
		t.Fatalf("activation after resend: %+v, %v", account, err)
	}
	mfa, err := authapplication.NewMFAService(repo, [32]byte{1})
	if err != nil {
		t.Fatal(err)
	}
	login, err := authapplication.NewLoginService(repo, repo, mfa)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := login.Login(ctx, domain.Login{Email: registration.Email, Password: registration.Password}); err != nil {
		t.Fatalf("original credentials no longer work: %v", err)
	}
}

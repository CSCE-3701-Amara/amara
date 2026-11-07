package database

import (
	"context"
	"crypto/sha256"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

func TestAccountActionTable(t *testing.T) {
	tests := []struct {
		accountType domain.AccountType
		table       string
		sessionType string
	}{
		{domain.AccountTypeCandidate, "candidates", "candidate"},
		{domain.AccountTypeInterviewer, "interviewers", "interviewer"},
		{domain.AccountTypeCompany, "companies", "company"},
	}

	for _, test := range tests {
		table, sessionType := accountActionTable(test.accountType)
		if table != test.table || sessionType != test.sessionType {
			t.Fatalf("accountActionTable(%q) = %q, %q; want %q, %q", test.accountType, table, sessionType, test.table, test.sessionType)
		}

	}

	if table, sessionType := accountActionTable("admin"); table != "" || sessionType != "" {
		t.Fatalf("accountActionTable(admin) = %q, %q; want empty values", table, sessionType)
	}
}

func TestActivateRequiresDatabase(t *testing.T) {
	if _, err := (*RegistrationRepository)(nil).Activate(context.Background(), "user@example.com", nil, time.Now()); err == nil {
		t.Fatal("expected missing database error")
	}
	repository := &RegistrationRepository{}
	if _, err := repository.Activate(context.Background(), "user@example.com", nil, time.Now()); err == nil {
		t.Fatal("expected missing database error")
	}
}

func TestActivateCandidateAccount(t *testing.T) {
	databaseURL := os.Getenv("AMARA_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("set AMARA_TEST_DATABASE_URL to run the PostgreSQL integration test")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create database pool: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping database: %v", err)
	}

	const accountID = "00000000-0000-0000-0000-000000000372"
	const email = "activation-test@example.com"
	const code = "12345678901234567890123456789012"
	hash := sha256.Sum256([]byte(code))
	_, _ = pool.Exec(ctx, "DELETE FROM auth_email_verification_deliveries WHERE account_id = $1::uuid", accountID)
	_, _ = pool.Exec(ctx, "DELETE FROM auth_email_verification_challenges WHERE account_id = $1::uuid", accountID)
	_, _ = pool.Exec(ctx, "DELETE FROM candidates WHERE id = $1::uuid", accountID)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM auth_email_verification_deliveries WHERE account_id = $1::uuid", accountID)
		_, _ = pool.Exec(context.Background(), "DELETE FROM auth_email_verification_challenges WHERE account_id = $1::uuid", accountID)
		_, _ = pool.Exec(context.Background(), "DELETE FROM candidates WHERE id = $1::uuid", accountID)
	})
	_, err = pool.Exec(ctx, `
			INSERT INTO candidates (id, name, email, password_hash, salt, created_at, account_status)
			VALUES ($1::uuid, 'Activation Test', $2, '', '', now(), 'pending_verification')`,
		accountID, email)
	if err != nil {
		t.Fatalf("insert candidate: %v", err)
	}
	_, err = pool.Exec(ctx, `
			INSERT INTO auth_email_verification_challenges
				(account_id, account_type, code_hash, expires_at)
			VALUES ($1::uuid, 'candidate', $2, now() + interval '1 hour')`,
		accountID, hash[:])
	if err != nil {
		t.Fatalf("insert challenge: %v", err)
	}
	_, err = pool.Exec(ctx, `
			INSERT INTO auth_email_verification_deliveries
				(account_id, account_type, recipient, expires_at)
			VALUES ($1::uuid, 'candidate', $2, now() + interval '1 hour')`,
		accountID, email)
	if err != nil {
		t.Fatalf("insert delivery: %v", err)
	}

	repository, err := NewRegistrationRepository(pool)
	if err != nil {
		t.Fatal(err)
	}
	account, err := repository.Activate(ctx, email, hash[:], time.Now().UTC())
	if err != nil {
		t.Fatalf("activate candidate: %v", err)
	}
	if account.ID != accountID || account.Type != domain.AccountTypeCandidate || account.Status != "active" {
		t.Fatalf("unexpected activated account: %+v", account)
	}
	if _, err := repository.Activate(ctx, email, hash[:], time.Now().UTC()); !errors.Is(err, domain.ErrInvalidVerificationCode) {
		t.Fatalf("expected consumed code to fail, got %v", err)
	}
}

func TestPerformAccountActionConcurrentSessions(t *testing.T) {
	databaseURL := os.Getenv("AMARA_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("set AMARA_TEST_DATABASE_URL to run the PostgreSQL integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create database pool: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping database: %v", err)
	}

	const accountID = "00000000-0000-0000-0000-000000000371"
	const firstToken = "account-action-test-session-1"
	const secondToken = "account-action-test-session-2"
	_, _ = pool.Exec(ctx, "DELETE FROM auth_sessions WHERE token IN ($1, $2)", firstToken, secondToken)
	_, _ = pool.Exec(ctx, "DELETE FROM companies WHERE id = $1::uuid", accountID)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM auth_sessions WHERE token IN ($1, $2)", firstToken, secondToken)
		_, _ = pool.Exec(context.Background(), "DELETE FROM companies WHERE id = $1::uuid", accountID)
	})

	_, err = pool.Exec(ctx, `
		INSERT INTO companies (id, email, password_hash, salt, legal_name, approved)
		VALUES ($1::uuid, $2, '', '', $3, TRUE)`,
		accountID, "account-action-test@example.com", "Account Action Test")
	if err != nil {
		t.Fatalf("insert test company: %v", err)
	}
	_, err = pool.Exec(ctx, `
		INSERT INTO auth_sessions (token, account_id, account_type, expires_at)
		VALUES ($1, $3::uuid, 'company', now() + interval '1 hour'),
		       ($2, $3::uuid, 'company', now() + interval '1 hour')`,
		firstToken, secondToken, accountID)
	if err != nil {
		t.Fatalf("insert test sessions: %v", err)
	}

	repository, err := NewRegistrationRepository(pool)
	if err != nil {
		t.Fatalf("create repository: %v", err)
	}

	errs := make(chan error, 2)
	var waitGroup sync.WaitGroup
	for _, token := range []string{firstToken, secondToken} {
		waitGroup.Add(1)
		go func(token string) {
			defer waitGroup.Done()
			errs <- repository.PerformAccountAction(ctx, accountID, token, domain.AccountTypeCompany, domain.AccountActionDelete)
		}(token)
	}
	waitGroup.Wait()
	close(errs)

	for err := range errs {
		if err != nil && !errors.Is(err, domain.ErrInvalidSession) {
			t.Fatalf("concurrent account action: %v", err)
		}
	}
}

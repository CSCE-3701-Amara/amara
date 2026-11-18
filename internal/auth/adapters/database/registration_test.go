package database

import (
	"context"
	"crypto/sha256"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
	"github.com/CSCE-3701-Amara/amara/internal/testutil"
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
	pool := testutil.Postgres(t)
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const accountID = "00000000-0000-0000-0000-000000000372"
	const email = "activation-test@example.com"
	const code = "12345678901234567890123456789012"
	const createdAt = "2026-01-02T03:04:05Z"
	hash := sha256.Sum256([]byte(code))
	_, err = pool.Exec(ctx, `
			INSERT INTO candidates (id, name, email, password_hash, salt, created_at, account_status)
			VALUES ($1::uuid, 'Activation Test', $2, '', '', $3::timestamp, 'pending_verification')`,
		accountID, email, createdAt)
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
	if account.CreatedAt != createdAt {
		t.Fatalf("activated account created_at = %q, want %q", account.CreatedAt, createdAt)
	}
	if _, err := repository.Activate(ctx, email, hash[:], time.Now().UTC()); !errors.Is(err, domain.ErrInvalidVerificationCode) {
		t.Fatalf("expected consumed code to fail, got %v", err)
	}
}

func TestActivateDistinguishesInvalidCodeFromDatabaseFailure(t *testing.T) {
	pool := testutil.Postgres(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	repository, err := NewRegistrationRepository(pool)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := repository.Activate(ctx, "unknown@example.com", []byte("unknown"), time.Now().UTC()); !errors.Is(err, domain.ErrInvalidVerificationCode) {
		t.Fatalf("unknown code error = %v, want invalid verification code", err)
	}

	pool.Close()
	if _, err := repository.Activate(ctx, "unknown@example.com", []byte("unknown"), time.Now().UTC()); err == nil || errors.Is(err, domain.ErrInvalidVerificationCode) {
		t.Fatalf("database failure error = %v, want non-validation error", err)
	}
}

func TestActivateInterviewerAccountReturnsCreationTimestamp(t *testing.T) {
	pool := testutil.Postgres(t)
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const accountID = "00000000-0000-0000-0000-000000000373"
	const companyID = "00000000-0000-0000-0000-000000000374"
	const email = "interviewer-activation-test@example.com"
	const code = "23456789012345678901234567890123"
	const createdAt = "2026-02-03T04:05:06Z"
	hash := sha256.Sum256([]byte(code))
	_, err = pool.Exec(ctx, `
		INSERT INTO companies (id, email, password_hash, salt, legal_name)
		VALUES ($1::uuid, $2, '', '', 'Activation Test Company')`,
		companyID, "activation-company@example.com")
	if err != nil {
		t.Fatalf("insert company: %v", err)
	}
	_, err = pool.Exec(ctx, `
		INSERT INTO interviewers (id, company_id, name, email, password_hash, salt, created_at, account_status)
		VALUES ($1::uuid, $2::uuid, 'Activation Test', $3, '', '', $4::timestamp, 'pending_verification')`,
		accountID, companyID, email, createdAt)
	if err != nil {
		t.Fatalf("insert interviewer: %v", err)
	}
	_, err = pool.Exec(ctx, `
		INSERT INTO auth_email_verification_challenges
			(account_id, account_type, code_hash, expires_at)
		VALUES ($1::uuid, 'interviewer', $2, now() + interval '1 hour')`,
		accountID, hash[:])
	if err != nil {
		t.Fatalf("insert challenge: %v", err)
	}
	_, err = pool.Exec(ctx, `
		INSERT INTO auth_email_verification_deliveries
			(account_id, account_type, recipient, expires_at)
		VALUES ($1::uuid, 'interviewer', $2, now() + interval '1 hour')`,
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
		t.Fatalf("activate interviewer: %v", err)
	}
	if account.ID != accountID || account.Type != domain.AccountTypeInterviewer || account.Status != "active" {
		t.Fatalf("unexpected activated account: %+v", account)
	}
	if account.CreatedAt != createdAt {
		t.Fatalf("activated account created_at = %q, want %q", account.CreatedAt, createdAt)
	}
}

func TestPerformAccountActionConcurrentSessions(t *testing.T) {
	pool := testutil.Postgres(t)
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const accountID = "00000000-0000-0000-0000-000000000371"
	const firstToken = "account-action-test-session-1"
	const secondToken = "account-action-test-session-2"

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
		sessionDigest(firstToken), sessionDigest(secondToken), accountID)
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

func TestBannedCandidateCannotUseOrCreateSessionPostgres(t *testing.T) {
	pool := testutil.Postgres(t)
	var err error
	ctx := t.Context()
	const id = "00000000-0000-0000-0000-000000000375"
	const token = "audit-banned-candidate"
	_, err = pool.Exec(ctx, `INSERT INTO candidates (id, name, email, password_hash, salt, created_at, account_status)
 VALUES ($1::uuid, 'Session test', 'audit-banned@example.test', '', '', now(), 'active')`, id)
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRegistrationRepository(pool)
	if err != nil {
		t.Fatal(err)
	}
	session := ports.Session{AccountID: id, AccountType: domain.AccountTypeCandidate, Token: token, ExpiresAt: time.Now().Add(time.Hour)}
	if err := repo.CreateSession(ctx, session); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.ValidateSession(ctx, token); err != nil {
		t.Fatalf("active candidate: %v", err)
	}
	if _, err := pool.Exec(ctx, "UPDATE candidates SET banned=true WHERE id=$1::uuid", id); err != nil {
		t.Fatal(err)
	}
	t.Run("existing session is denied after ban", func(t *testing.T) {
		if _, err := repo.ValidateSession(ctx, token); !errors.Is(err, domain.ErrInvalidSession) {
			t.Fatalf("banned session validation = %v, want invalid session", err)
		}
	})
	t.Run("ban racing with login prevents session creation", func(t *testing.T) {
		session.Token = token + "-new"
		if err := repo.CreateSession(ctx, session); !errors.Is(err, domain.ErrAccountInactive) {
			t.Fatalf("create banned session = %v, want inactive", err)
		}
	})
}

// The table holds token digests only: a stored row must not work as a cookie,
// and the migration must convert existing plaintext sessions in place.
func TestSessionsStoreTokenDigestsPostgres(t *testing.T) {
	pool := testutil.Postgres(t)
	ctx := t.Context()
	const id = "00000000-0000-0000-0000-000000000376"
	if _, err := pool.Exec(ctx, `INSERT INTO candidates (id, name, email, password_hash, salt, created_at, account_status)
 VALUES ($1::uuid, 'Digest test', 'digest@example.test', '', '', now(), 'active')`, id); err != nil {
		t.Fatal(err)
	}
	repo, err := NewRegistrationRepository(pool)
	if err != nil {
		t.Fatal(err)
	}
	const token = "digest-test-token-0123456789abcdefghijklmnopq"
	session := ports.Session{Token: token, AccountID: id, AccountType: domain.AccountTypeCandidate, ExpiresAt: time.Now().Add(time.Hour)}
	if err := repo.CreateSession(ctx, session); err != nil {
		t.Fatal(err)
	}
	var stored string
	if err := pool.QueryRow(ctx, "SELECT token FROM auth_sessions WHERE account_id=$1::uuid", id).Scan(&stored); err != nil {
		t.Fatal(err)
	}
	if stored == token || stored != sessionDigest(token) {
		t.Fatalf("stored %q, want the digest of the token", stored)
	}
	if _, err := repo.ValidateSession(ctx, token); err != nil {
		t.Fatalf("token no longer validates: %v", err)
	}
	if _, err := repo.ValidateSession(ctx, stored); !errors.Is(err, domain.ErrInvalidSession) {
		t.Fatalf("a stored digest was accepted as a token: %v", err)
	}

	// Legacy plaintext row, then the migration's conversion statement.
	const legacy = "legacy-plaintext-session-token-abcdefghijklm"
	if _, err := pool.Exec(ctx, `INSERT INTO auth_sessions (token, account_id, account_type, expires_at)
 VALUES ($1, $2::uuid, 'candidate', now() + interval '1 hour')`, legacy, id); err != nil {
		t.Fatal(err)
	}
	migration, err := os.ReadFile("../../../../supabase/migrations/20261118000000_hash_auth_session_tokens.sql")
	if err != nil {
		t.Fatal(err)
	}
	for range 2 { // idempotent
		if _, err := pool.Exec(ctx, string(migration)); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := repo.ValidateSession(ctx, legacy); err != nil {
		t.Fatalf("existing session was signed out by the migration: %v", err)
	}
	if _, err := repo.ValidateSession(ctx, token); err != nil {
		t.Fatalf("migration re-hashed an already hashed session: %v", err)
	}
}

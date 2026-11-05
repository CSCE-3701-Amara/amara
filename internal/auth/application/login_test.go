package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/argon2"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type loginCredentials struct {
	record ports.Credentials
	err    error
}

func (f loginCredentials) FindCredentials(context.Context, string) (ports.Credentials, error) {
	return f.record, f.err
}

type loginSessions struct {
	session ports.Session
	err     error
}

func (f *loginSessions) CreateSession(_ context.Context, session ports.Session) error {
	f.session = session
	return f.err
}

type loginMFA struct {
	record ports.MFARecord
	err    error
}

func (f loginMFA) GetMFASecret(context.Context, string, domain.AccountType) (ports.MFARecord, error) {
	return f.record, f.err
}
func (f loginMFA) SaveMFASecret(context.Context, ports.MFARecord) error        { return nil }
func (f loginMFA) EnableMFA(context.Context, string, domain.AccountType) error { return nil }

func loginService(t *testing.T, credentials ports.Credentials, mfa loginMFA, sessions *loginSessions) *LoginService {
	t.Helper()
	mfaService, err := NewMFAService(mfa, [32]byte{1})
	if err != nil {
		t.Fatal(err)
	}
	service, err := NewLoginService(loginCredentials{record: credentials}, sessions, mfaService)
	if err != nil {
		t.Fatal(err)
	}
	service.now = func() time.Time { return time.Unix(1000, 0) }
	return service
}

func hashedPassword(password string) ([]byte, []byte) {
	salt := []byte("01234567890123456789")
	return argon2.IDKey([]byte(password), salt, 3, 64*1024, 1, 32), salt
}

func TestLoginCreatesSession(t *testing.T) {
	hash, salt := hashedPassword("password123")
	sessions := &loginSessions{}
	service := loginService(t, ports.Credentials{
		ID: "account", Email: "User@Example.com", AccountType: domain.AccountTypeCandidate,
		PasswordHash: hash, Salt: salt, Status: "active",
	}, loginMFA{err: ports.ErrMFANotConfigured}, sessions)

	session, err := service.Login(context.Background(), domain.Login{
		Email: " user@example.com ", Password: "password123",
	})
	if err != nil {
		t.Fatal(err)
	}
	if session.Email != "User@Example.com" || sessions.session.AccountID != "account" ||
		sessions.session.Token == "" || session.ExpiresAt != "1970-01-02T00:16:40Z" {
		t.Fatalf("unexpected session: %+v", session)
	}
}

func TestLoginRejectsInvalidAndInactiveAccounts(t *testing.T) {
	hash, salt := hashedPassword("password123")
	tests := []struct {
		name     string
		record   ports.Credentials
		password string
		want     error
	}{
		{"wrong password", ports.Credentials{PasswordHash: hash, Salt: salt, Status: "active"}, "wrong", domain.ErrInvalidCredentials},
		{"banned", ports.Credentials{PasswordHash: hash, Salt: salt, Status: "active", Banned: true}, "password123", domain.ErrAccountForbidden},
		{"pending", ports.Credentials{PasswordHash: hash, Salt: salt, Status: "pending_verification"}, "password123", domain.ErrAccountInactive},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := loginService(t, test.record, loginMFA{err: ports.ErrMFANotConfigured}, &loginSessions{})
			_, err := service.Login(context.Background(), domain.Login{Email: "user@example.com", Password: test.password})
			if !errors.Is(err, test.want) {
				t.Fatalf("got %v, want %v", err, test.want)
			}
		})
	}
}

func TestLoginRequiresMFAForInterviewer(t *testing.T) {
	hash, salt := hashedPassword("password123")
	service := loginService(t, ports.Credentials{
		ID: "account", AccountType: domain.AccountTypeInterviewer, PasswordHash: hash, Salt: salt,
		Status: "active", CompanyApproved: true,
	}, loginMFA{err: ports.ErrMFANotConfigured}, &loginSessions{})
	_, err := service.Login(context.Background(), domain.Login{Email: "user@example.com", Password: "password123"})
	if !errors.Is(err, domain.ErrMFARequired) {
		t.Fatalf("got %v, want MFA required", err)
	}
}

func TestLoginDependencyAndPersistenceErrors(t *testing.T) {
	if _, err := NewLoginService(nil, &loginSessions{}, &MFAService{}); err == nil {
		t.Fatal("expected dependency validation error")
	}
	hash, salt := hashedPassword("password123")
	sessions := &loginSessions{err: errors.New("storage down")}
	service := loginService(t, ports.Credentials{
		ID: "account", AccountType: domain.AccountTypeCandidate, PasswordHash: hash, Salt: salt, Status: "active",
	}, loginMFA{err: ports.ErrMFANotConfigured}, sessions)
	if _, err := service.Login(context.Background(), domain.Login{Email: "user@example.com", Password: "password123"}); err == nil {
		t.Fatal("expected session persistence error")
	}
	service = loginService(t, ports.Credentials{}, loginMFA{err: ports.ErrMFANotConfigured}, &loginSessions{})
	service.credentials = loginCredentials{err: errors.New("database down")}
	if _, err := service.Login(context.Background(), domain.Login{Email: "user@example.com", Password: "password123"}); err == nil {
		t.Fatal("expected credential lookup error")
	}
}

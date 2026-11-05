package application

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

const sessionLifetime = 24 * time.Hour

type LoginService struct {
	credentials ports.CredentialRepository
	sessions    ports.SessionRepository
	mfa         *MFAService
	now         func() time.Time
}

func NewLoginService(credentials ports.CredentialRepository, sessions ports.SessionRepository, mfa *MFAService) (*LoginService, error) {
	if credentials == nil || sessions == nil || mfa == nil {
		return nil, errors.New("login dependencies must not be nil")
	}
	return &LoginService{credentials: credentials, sessions: sessions, mfa: mfa, now: time.Now}, nil
}

func (s *LoginService) Login(ctx context.Context, login domain.Login) (domain.Session, error) {
	login.Email = strings.ToLower(strings.TrimSpace(login.Email))
	if login.Email == "" || login.Password == "" {
		return domain.Session{}, domain.ErrInvalidCredentials
	}
	credentials, err := s.credentials.FindCredentials(ctx, login.Email)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return domain.Session{}, err
		}
		return domain.Session{}, fmt.Errorf("find login credentials: %w", err)
	}
	if !verifyPassword(login.Password, credentials.PasswordHash, credentials.Salt) {
		return domain.Session{}, domain.ErrInvalidCredentials
	}
	if credentials.Banned {
		return domain.Session{}, domain.ErrAccountForbidden
	}
	if credentials.Deactivated || credentials.Status != "active" {
		return domain.Session{}, domain.ErrAccountInactive
	}
	if credentials.AccountType == domain.AccountTypeInterviewer && !credentials.CompanyApproved {
		return domain.Session{}, domain.ErrAccountForbidden
	}
	if err := s.mfa.VerifyLoginCode(ctx, credentials.ID, credentials.AccountType, login.MFACode); err != nil {
		return domain.Session{}, err
	}

	expiresAt := s.now().UTC().Add(sessionLifetime)
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return domain.Session{}, fmt.Errorf("generate session token: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	if err := s.sessions.CreateSession(ctx, ports.Session{
		Token: token, AccountID: credentials.ID, AccountType: credentials.AccountType, ExpiresAt: expiresAt,
	}); err != nil {
		return domain.Session{}, fmt.Errorf("create session: %w", err)
	}
	return domain.Session{
		Token: token, AccountID: credentials.ID, AccountType: credentials.AccountType,
		Email: credentials.Email, ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}

func verifyPassword(password string, expected, salt []byte) bool {
	actual := argon2.IDKey([]byte(password), salt, 3, 64*1024, 1, 32)
	return len(actual) == len(expected) && subtle.ConstantTimeCompare(actual, expected) == 1
}

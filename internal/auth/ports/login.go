package ports

import (
	"context"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

type CredentialRepository interface {
	FindCredentials(context.Context, string) (Credentials, error)
}

type Credentials struct {
	ID              string
	Email           string
	AccountType     domain.AccountType
	PasswordHash    []byte
	Salt            []byte
	Status          string
	Deactivated     bool
	Banned          bool
	CompanyID       string
	CompanyApproved bool
}

type SessionRepository interface {
	CreateSession(context.Context, Session) error
}

type Session struct {
	Token       string
	AccountID   string
	AccountType domain.AccountType
	ExpiresAt   time.Time
}

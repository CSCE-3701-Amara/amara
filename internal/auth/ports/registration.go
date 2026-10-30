package ports

import (
	"context"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

type AccountRepository interface {
	EmailExists(context.Context, string) (bool, error)
	FindApprovedCompanyBySecret(context.Context, string) (string, error)
	CreateCandidate(context.Context, string, string, []byte, []byte, time.Time) (domain.Account, error)
	CreateInterviewer(context.Context, string, string, []byte, []byte, string, time.Time) (domain.Account, error)
	SaveVerificationChallenge(context.Context, VerificationChallenge) error
}

type PasswordCompromiseChecker interface {
	IsCompromised(context.Context, string) (bool, error)
}

type VerificationSender interface {
	SendVerificationCode(context.Context, string, string) error
}

type VerificationChallenge struct {
	AccountID   string
	AccountType domain.AccountType
	CodeHash    []byte
	ExpiresAt   time.Time
}

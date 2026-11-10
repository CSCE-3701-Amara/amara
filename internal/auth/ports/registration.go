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
	FindPendingVerification(context.Context, string) (PendingVerification, error)
}

type CompanyRegistrationRepository interface {
	EmailExists(context.Context, string) (bool, error)
	RegistrationNumberExists(context.Context, string) (bool, error)
	CreateCompanyRegistration(context.Context, CompanyRegistrationRecord) (domain.Account, error)
}

type CompanyRegistrationRecord struct {
	AutoApprove        bool
	CompanyName        string
	RegistrationNumber string
	Email              string
	PasswordHash       []byte
	Salt               []byte
	CreatedAt          time.Time
}

type PasswordCompromiseChecker interface {
	IsCompromised(context.Context, string) (bool, error)
}

type VerificationSender interface {
	SendVerificationCode(context.Context, string, string) error
}

type ActivationRepository interface {
	Activate(context.Context, string, []byte, time.Time) (domain.Account, error)
}

type VerificationChallenge struct {
	AccountID   string
	AccountType domain.AccountType
	Recipient   string
	CodeHash    []byte
	ExpiresAt   time.Time
}

type PendingVerification struct {
	PasswordHash []byte
	Salt         []byte
	Account      domain.Account
	Recipient    string
	ExpiresAt    time.Time
}

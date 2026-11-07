package ports

import (
	"context"
	"errors"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

var ErrMFANotConfigured = errors.New("MFA is not configured")

type MFARepository interface {
	GetMFASecret(context.Context, string, domain.AccountType) (MFARecord, error)
	SaveMFASecret(context.Context, MFARecord) error
	EnableMFA(context.Context, string, domain.AccountType) error
}

type MFARecord struct {
	AccountID   string
	AccountType domain.AccountType
	Secret      []byte
	Enabled     bool
	CreatedAt   time.Time
}

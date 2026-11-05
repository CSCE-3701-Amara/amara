package application

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/argon2"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

const verificationCodeLifetime = 15 * time.Minute

type RegistrationService struct {
	repository ports.AccountRepository
	checker    ports.PasswordCompromiseChecker
	sender     ports.VerificationSender
	now        func() time.Time
}

func NewRegistrationService(repository ports.AccountRepository, checker ports.PasswordCompromiseChecker, sender ports.VerificationSender) (*RegistrationService, error) {
	if repository == nil || checker == nil || sender == nil {
		return nil, errors.New("registration dependencies must not be nil")
	}
	return &RegistrationService{repository: repository, checker: checker, sender: sender, now: time.Now}, nil
}

func (s *RegistrationService) Register(ctx context.Context, registration domain.Registration) (domain.Account, error) {
	registration = registration.Normalized()
	if err := registration.Validate(); err != nil {
		return domain.Account{}, err
	}

	exists, err := s.repository.EmailExists(ctx, registration.Email)
	if err != nil {
		return domain.Account{}, fmt.Errorf("check registration email: %w", err)
	}
	if exists {
		return domain.Account{}, domain.ErrEmailAlreadyExists
	}

	compromised, err := s.checker.IsCompromised(ctx, registration.Password)
	if err != nil {
		return domain.Account{}, fmt.Errorf("check password compromise: %w", err)
	}
	if compromised {
		return domain.Account{}, fmt.Errorf("%w: password has appeared in a data breach", domain.ErrInvalidRegistration)
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return domain.Account{}, fmt.Errorf("generate password salt: %w", err)
	}
	passwordHash := argon2.IDKey([]byte(registration.Password), salt, 3, 64*1024, 1, 32)
	createdAt := s.now().UTC()
	var account domain.Account

	switch registration.AccountType {
	case domain.AccountTypeCandidate:
		account, err = s.repository.CreateCandidate(ctx, registration.FirstName+" "+registration.LastName, registration.Email, passwordHash, salt, createdAt)
	case domain.AccountTypeInterviewer:
		companyID, resolveErr := s.repository.FindApprovedCompanyBySecret(ctx, registration.CompanySecret)
		if resolveErr != nil {
			if errors.Is(resolveErr, domain.ErrCompanyNotFound) {
				return domain.Account{}, resolveErr
			}
			return domain.Account{}, fmt.Errorf("resolve interviewer company: %w", resolveErr)
		}
		account, err = s.repository.CreateInterviewer(ctx, registration.FirstName+" "+registration.LastName, registration.Email, passwordHash, salt, companyID, createdAt)
	}
	if err != nil {
		return domain.Account{}, fmt.Errorf("create %s account: %w", registration.AccountType, err)
	}

	code := make([]byte, 32)
	if _, err := rand.Read(code); err != nil {
		return domain.Account{}, fmt.Errorf("generate verification code: %w", err)
	}
	codeHash := sha256.Sum256(code)
	challenge := ports.VerificationChallenge{
		AccountID:   account.ID,
		AccountType: registration.AccountType,
		CodeHash:    codeHash[:],
		ExpiresAt:   createdAt.Add(verificationCodeLifetime),
	}
	if err := s.repository.SaveVerificationChallenge(ctx, challenge); err != nil {
		return domain.Account{}, fmt.Errorf("save verification challenge: %w", err)
	}
	if err := s.sender.SendVerificationCode(ctx, registration.Email, fmt.Sprintf("%x", code)); err != nil {
		return domain.Account{}, fmt.Errorf("send verification code: %w", err)
	}
	return account, nil
}

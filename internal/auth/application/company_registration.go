package application

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/argon2"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type CompanyRegistrationService struct {
	repository ports.CompanyRegistrationRepository
	checker    ports.PasswordCompromiseChecker
	now        func() time.Time
}

func NewCompanyRegistrationService(repository ports.CompanyRegistrationRepository, checker ports.PasswordCompromiseChecker) (*CompanyRegistrationService, error) {
	if repository == nil || checker == nil {
		return nil, errors.New("company registration dependencies must not be nil")
	}
	return &CompanyRegistrationService{repository: repository, checker: checker, now: time.Now}, nil
}

func (s *CompanyRegistrationService) Register(ctx context.Context, registration domain.CompanyRegistration) (domain.Account, error) {
	registration = registration.Normalized()
	if err := registration.Validate(); err != nil {
		return domain.Account{}, err
	}

	emailExists, err := s.repository.EmailExists(ctx, registration.Email)
	if err != nil {
		return domain.Account{}, fmt.Errorf("check company registration email: %w", err)
	}
	if emailExists {
		return domain.Account{}, domain.ErrEmailAlreadyExists
	}
	registrationNumberExists, err := s.repository.RegistrationNumberExists(ctx, registration.RegistrationNumber)
	if err != nil {
		return domain.Account{}, fmt.Errorf("check registration number: %w", err)
	}
	if registrationNumberExists {
		return domain.Account{}, domain.ErrCompanyRegistrationExists
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
	account, err := s.repository.CreateCompanyRegistration(ctx, ports.CompanyRegistrationRecord{
		CompanyName: registration.CompanyName, RegistrationNumber: registration.RegistrationNumber,
		Email:        registration.Email,
		PasswordHash: argon2.IDKey([]byte(registration.Password), salt, 3, 64*1024, 1, 32),
		Salt:         salt, CreatedAt: s.now().UTC(),
	})
	if err != nil {
		return domain.Account{}, fmt.Errorf("create company registration: %w", err)
	}
	return account, nil
}

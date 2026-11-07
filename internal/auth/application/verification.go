package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type VerificationService struct {
	repository ports.ActivationRepository
	now        func() time.Time
}

func NewVerificationService(repository ports.ActivationRepository) (*VerificationService, error) {
	if repository == nil {
		return nil, errors.New("verification repository must not be nil")
	}
	return &VerificationService{repository: repository, now: time.Now}, nil
}

func (s *VerificationService) Verify(ctx context.Context, email, code string) (domain.Account, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	code = strings.TrimSpace(code)
	if email == "" || len(code) != 64 {
		return domain.Account{}, domain.ErrInvalidVerificationCode
	}
	rawCode, err := hex.DecodeString(code)
	if err != nil {
		return domain.Account{}, domain.ErrInvalidVerificationCode
	}
	codeHash := sha256.Sum256(rawCode)
	account, err := s.repository.Activate(ctx, email, codeHash[:], s.now().UTC())
	if err != nil {
		if errors.Is(err, domain.ErrInvalidVerificationCode) {
			return domain.Account{}, err
		}
		return domain.Account{}, fmt.Errorf("activate account: %w", err)
	}
	return account, nil
}

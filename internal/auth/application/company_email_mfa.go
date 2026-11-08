package application

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

type CompanyEmailChallengeRepository interface {
	ReserveCompanyEmailCode(context.Context, string, []byte, time.Time, time.Time) (bool, error)
	ReleaseCompanyEmailCode(context.Context, string, []byte) error
	ConsumeCompanyEmailCode(context.Context, string, []byte) (bool, error)
}

type LoginCodeSender interface {
	SendLoginCode(context.Context, string, string) error
}

type companyEmailMFA struct {
	repository CompanyEmailChallengeRepository
	sender     LoginCodeSender
}

// ConfigureCompanyEmail enables email challenges for companies without enrolled TOTP.
func (s *MFAService) ConfigureCompanyEmail(repository CompanyEmailChallengeRepository, sender LoginCodeSender) {
	s.companyEmail = &companyEmailMFA{repository: repository, sender: sender}
}

var emailLoginCode = regexp.MustCompile(`^[0-9]{6}$`)

func (s *MFAService) verifyLogin(ctx context.Context, accountID string, accountType domain.AccountType, email, code string) error {
	err := s.VerifyLoginCode(ctx, accountID, accountType, code)
	if accountType != domain.AccountTypeCompany || !errors.Is(err, domain.ErrMFARequired) || s.companyEmail == nil {
		return err
	}
	code = strings.TrimSpace(code)
	if code != "" {
		if !emailLoginCode.MatchString(code) {
			return domain.ErrInvalidMFACode
		}
		ok, err := s.companyEmail.repository.ConsumeCompanyEmailCode(ctx, accountID, s.loginCodeMAC(accountID, code))
		if err != nil {
			return fmt.Errorf("verify company email code: %w", err)
		}
		if !ok {
			return domain.ErrInvalidMFACode
		}
		return nil
	}

	number, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return fmt.Errorf("generate company email code: %w", err)
	}
	code = fmt.Sprintf("%06d", number.Int64())
	now := s.now()
	mac := s.loginCodeMAC(accountID, code)
	reserved, err := s.companyEmail.repository.ReserveCompanyEmailCode(ctx, accountID, mac, now.Add(10*time.Minute), now.Add(time.Minute))
	if err != nil {
		return fmt.Errorf("reserve company email code: %w", err)
	}
	if reserved {
		if err := s.companyEmail.sender.SendLoginCode(ctx, email, code); err != nil {
			cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
			defer cancel()
			releaseErr := s.companyEmail.repository.ReleaseCompanyEmailCode(cleanup, accountID, mac)
			return fmt.Errorf("send company email code: %w", errors.Join(err, releaseErr))
		}
	}
	return domain.ErrEmailMFARequired
}

func (s *MFAService) loginCodeMAC(accountID, code string) []byte {
	mac := hmac.New(sha256.New, s.key[:])
	_, _ = mac.Write([]byte("company-email-mfa:" + accountID + ":" + code))
	return mac.Sum(nil)
}

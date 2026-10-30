package application

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

const (
	totpPeriod = 30 * time.Second
	totpDigits = 6
)

type MFAService struct {
	repository ports.MFARepository
	key        [32]byte
	now        func() time.Time
}

func NewMFAService(repository ports.MFARepository, encryptionKey [32]byte) (*MFAService, error) {
	if repository == nil {
		return nil, errors.New("MFA repository must not be nil")
	}
	if encryptionKey == [32]byte{} {
		return nil, errors.New("MFA encryption key must not be empty")
	}
	return &MFAService{repository: repository, key: encryptionKey, now: time.Now}, nil
}

func (s *MFAService) BeginEnrollment(ctx context.Context, accountID string, accountType domain.AccountType, email string) (domain.MFAEnrollment, error) {
	record, err := s.repository.GetMFASecret(ctx, accountID, accountType)
	if err == nil && record.Enabled {
		return domain.MFAEnrollment{}, domain.ErrMFAAlreadyEnabled
	}
	if err != nil && !errors.Is(err, ports.ErrMFANotConfigured) {
		return domain.MFAEnrollment{}, fmt.Errorf("load MFA enrollment: %w", err)
	}

	secret := make([]byte, 20)
	if _, err := rand.Read(secret); err != nil {
		return domain.MFAEnrollment{}, fmt.Errorf("generate MFA secret: %w", err)
	}
	encrypted, err := s.encrypt(secret)
	if err != nil {
		return domain.MFAEnrollment{}, err
	}
	if err := s.repository.SaveMFASecret(ctx, ports.MFARecord{
		AccountID: accountID, AccountType: accountType, Secret: encrypted, CreatedAt: s.now().UTC(),
	}); err != nil {
		return domain.MFAEnrollment{}, fmt.Errorf("save MFA enrollment: %w", err)
	}
	encoded := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret)
	return domain.MFAEnrollment{
		AccountID: accountID, AccountType: accountType, Email: email,
		ProvisioningURI: provisioningURI(email, accountType, encoded),
	}, nil
}

func (s *MFAService) ConfirmEnrollment(ctx context.Context, accountID string, accountType domain.AccountType, code string) error {
	record, err := s.repository.GetMFASecret(ctx, accountID, accountType)
	if err != nil {
		return fmt.Errorf("load MFA enrollment: %w", err)
	}
	secret, err := s.decrypt(record.Secret)
	if err != nil {
		return fmt.Errorf("decrypt MFA secret: %w", err)
	}
	if !verifyTOTP(secret, code, s.now()) {
		return domain.ErrInvalidMFACode
	}
	if err := s.repository.EnableMFA(ctx, accountID, accountType); err != nil {
		return fmt.Errorf("enable MFA: %w", err)
	}
	return nil
}

func (s *MFAService) VerifyLoginCode(ctx context.Context, accountID string, accountType domain.AccountType, code string) error {
	record, err := s.repository.GetMFASecret(ctx, accountID, accountType)
	if err != nil {
		if errors.Is(err, ports.ErrMFANotConfigured) && domain.MFARequired(accountType) {
			return domain.ErrMFARequired
		}
		return fmt.Errorf("load MFA credentials: %w", err)
	}
	if !record.Enabled {
		if domain.MFARequired(accountType) {
			return domain.ErrMFARequired
		}
		return nil
	}
	secret, err := s.decrypt(record.Secret)
	if err != nil {
		return fmt.Errorf("decrypt MFA secret: %w", err)
	}
	if !verifyTOTP(secret, code, s.now()) {
		return domain.ErrInvalidMFACode
	}
	return nil
}

func (s *MFAService) encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key[:])
	if err != nil {
		return nil, fmt.Errorf("create MFA cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create MFA cipher mode: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("generate MFA nonce: %w", err)
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (s *MFAService) decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("invalid encrypted MFA secret")
	}
	nonce, payload := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, payload, nil)
}

func provisioningURI(email string, accountType domain.AccountType, secret string) string {
	label := url.PathEscape("AMARA:" + email)
	query := url.Values{}
	query.Set("secret", secret)
	query.Set("issuer", "AMARA")
	query.Set("algorithm", "SHA1")
	query.Set("digits", strconv.Itoa(totpDigits))
	query.Set("period", strconv.Itoa(int(totpPeriod/time.Second)))
	return "otpauth://totp/" + label + "?" + query.Encode()
}

func verifyTOTP(secret []byte, code string, now time.Time) bool {
	if len(code) != totpDigits {
		return false
	}
	if _, err := strconv.Atoi(code); err != nil {
		return false
	}
	current := now.Unix() / int64(totpPeriod/time.Second)
	for offset := int64(-1); offset <= 1; offset++ {
		if hmac.Equal([]byte(generateTOTP(secret, current+offset)), []byte(code)) {
			return true
		}
	}
	return false
}

func generateTOTP(secret []byte, counter int64) string {
	var input [8]byte
	binary.BigEndian.PutUint64(input[:], uint64(counter))
	mac := hmac.New(sha1.New, secret)
	_, _ = mac.Write(input[:])
	sum := mac.Sum(nil)
	offset := sum[len(sum)-1] & 0x0f
	value := (uint32(sum[offset])&0x7f)<<24 |
		uint32(sum[offset+1])<<16 |
		uint32(sum[offset+2])<<8 |
		uint32(sum[offset+3])
	return fmt.Sprintf("%06d", value%1000000)
}

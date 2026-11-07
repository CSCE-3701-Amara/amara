package application

import (
	"context"
	"errors"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type mfaRepository struct {
	record    ports.MFARecord
	getErr    error
	saveErr   error
	enableErr error
	saved     ports.MFARecord
	enabled   bool
}

func (r *mfaRepository) GetMFASecret(context.Context, string, domain.AccountType) (ports.MFARecord, error) {
	if r.getErr != nil {
		return ports.MFARecord{}, r.getErr
	}
	if r.record.AccountID == "" {
		return ports.MFARecord{}, ports.ErrMFANotConfigured
	}
	return r.record, nil
}
func (r *mfaRepository) SaveMFASecret(_ context.Context, record ports.MFARecord) error {
	r.saved = record
	r.record = record
	return r.saveErr
}
func (r *mfaRepository) EnableMFA(context.Context, string, domain.AccountType) error {
	r.enabled = true
	return r.enableErr
}

func newMFAService(t *testing.T, repository *mfaRepository) *MFAService {
	t.Helper()
	var key [32]byte
	for index := range key {
		key[index] = byte(index + 1)
	}
	service, err := NewMFAService(repository, key)
	if err != nil {
		t.Fatal(err)
	}
	service.now = func() time.Time { return time.Unix(1_700_000_000, 0).UTC() }
	return service
}

func TestNewMFAServiceRequiresDependencies(t *testing.T) {
	var key [32]byte
	if _, err := NewMFAService(nil, key); err == nil {
		t.Fatal("expected repository error")
	}
	if _, err := NewMFAService(&mfaRepository{}, key); err == nil {
		t.Fatal("expected encryption key error")
	}
}

func TestMFAEnrollmentAndVerification(t *testing.T) {
	repository := &mfaRepository{}
	service := newMFAService(t, repository)
	enrollment, err := service.BeginEnrollment(context.Background(), "account-id", domain.AccountTypeCandidate, "user@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(enrollment.ProvisioningURI, "otpauth://totp/") || len(repository.saved.Secret) <= 20 {
		t.Fatalf("unexpected enrollment: %+v", enrollment)
	}
	secret, err := service.decrypt(repository.saved.Secret)
	if err != nil {
		t.Fatal(err)
	}
	code := generateTOTP(secret, service.now().Unix()/30)
	if err := service.ConfirmEnrollment(context.Background(), "account-id", domain.AccountTypeCandidate, code); err != nil {
		t.Fatal(err)
	}
	if !repository.enabled {
		t.Fatal("expected MFA to be enabled")
	}
	repository.record.Enabled = true
	if err := service.VerifyLoginCode(context.Background(), "account-id", domain.AccountTypeCandidate, code); err != nil {
		t.Fatal(err)
	}
}

func TestMFARequiredAndInvalidCodes(t *testing.T) {
	service := newMFAService(t, &mfaRepository{})
	if err := service.VerifyLoginCode(context.Background(), "id", domain.AccountTypeInterviewer, "123456"); !errors.Is(err, domain.ErrMFARequired) {
		t.Fatalf("expected MFA required, got %v", err)
	}
	repository := &mfaRepository{}
	service = newMFAService(t, repository)
	if _, err := service.BeginEnrollment(context.Background(), "id", domain.AccountTypeInterviewer, "user@example.com"); err != nil {
		t.Fatal(err)
	}
	repository.record.Enabled = true
	if err := service.VerifyLoginCode(context.Background(), "id", domain.AccountTypeInterviewer, "123"); !errors.Is(err, domain.ErrInvalidMFACode) {
		t.Fatalf("expected invalid MFA code, got %v", err)
	}
}

func TestMFAEnrollmentErrors(t *testing.T) {
	repository := &mfaRepository{record: ports.MFARecord{AccountID: "id", Enabled: true}}
	if err := func() error {
		service := newMFAService(t, repository)
		_, err := service.BeginEnrollment(context.Background(), "id", domain.AccountTypeCandidate, "user@example.com")
		return err
	}(); !errors.Is(err, domain.ErrMFAAlreadyEnabled) {
		t.Fatalf("expected already enabled, got %v", err)
	}
	repository = &mfaRepository{getErr: errors.New("db down")}
	if _, err := newMFAService(t, repository).BeginEnrollment(context.Background(), "id", domain.AccountTypeCandidate, "user@example.com"); err == nil {
		t.Fatal("expected load error")
	}
	repository = &mfaRepository{saveErr: errors.New("db down")}
	if _, err := newMFAService(t, repository).BeginEnrollment(context.Background(), "id", domain.AccountTypeCandidate, "user@example.com"); err == nil {
		t.Fatal("expected save error")
	}
}

func TestGenerateTOTPRFC4226Vectors(t *testing.T) {
	// RFC 4226 Appendix D: independent reference values for the HOTP primitive
	// used by TOTP. Keep the six-digit output and SHA-1 enrollment compatible.
	secret := []byte("12345678901234567890")
	for counter, want := range []string{
		"755224", "287082", "359152", "969429", "338314",
		"254676", "287922", "162583", "399871", "520489",
	} {
		if got := generateTOTP(secret, int64(counter)); got != want {
			t.Fatalf("counter %d: code = %q, want %q", counter, got, want)
		}
	}
}

func TestGenerateTOTPRejectsNegativeCounters(t *testing.T) {
	for _, counter := range []int64{math.MinInt64, -1} {
		if got := generateTOTP([]byte("12345678901234567890"), counter); got != "" {
			t.Fatalf("negative counter %d must not produce a code", counter)
		}
	}
}

func TestVerifyTOTPEpochAndWindowBoundaries(t *testing.T) {
	secret := []byte("12345678901234567890")
	for _, tc := range []struct {
		name    string
		seconds int64
		code    string
		want    bool
	}{
		{"before epoch", -1, "755224", false},
		{"negative time step", -30, "755224", false},
		{"epoch", 0, "755224", true},
		{"next epoch step", 0, "287082", true},
		{"previous step", 60, "287082", true},
		{"current step", 60, "359152", true},
		{"next step", 60, "969429", true},
		{"expired step", 60, "755224", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := verifyTOTP(secret, tc.code, time.Unix(tc.seconds, 0)); got != tc.want {
				t.Fatalf("verification = %t, want %t", got, tc.want)
			}
		})
	}
}

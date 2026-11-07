package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

type activationRepository struct {
	email string
	hash  []byte
	now   time.Time
	err   error
}

func (r *activationRepository) Activate(_ context.Context, email string, hash []byte, now time.Time) (domain.Account, error) {
	if r.err != nil {
		return domain.Account{}, r.err
	}
	r.email, r.hash, r.now = email, hash, now
	return domain.Account{ID: "candidate-id", Type: domain.AccountTypeCandidate, Email: email, Status: "active"}, nil
}

func TestVerifyActivatesAccount(t *testing.T) {
	repository := &activationRepository{}
	service, err := NewVerificationService(repository)
	if err != nil {
		t.Fatal(err)
	}
	service.now = func() time.Time { return time.Date(2026, 11, 7, 18, 0, 0, 0, time.UTC) }
	rawCode := []byte("12345678901234567890123456789012")
	code := hex.EncodeToString(rawCode)
	hash := sha256.Sum256(rawCode)
	account, err := service.Verify(context.Background(), " User@Example.com ", code)
	if err != nil {
		t.Fatal(err)
	}
	if account.Status != "active" || repository.email != "user@example.com" || string(repository.hash) != string(hash[:]) {
		t.Fatalf("unexpected activation: account=%+v email=%q hash=%x", account, repository.email, repository.hash)
	}
}

func TestVerifyRejectsMalformedCode(t *testing.T) {
	service, err := NewVerificationService(&activationRepository{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Verify(context.Background(), "user@example.com", "not-a-code"); !errors.Is(err, domain.ErrInvalidVerificationCode) {
		t.Fatalf("expected invalid code, got %v", err)
	}
}

func TestNewVerificationServiceRequiresRepository(t *testing.T) {
	if _, err := NewVerificationService(nil); err == nil {
		t.Fatal("expected missing repository error")
	}
}

func TestVerifyRejectsInvalidHexCode(t *testing.T) {
	service, err := NewVerificationService(&activationRepository{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Verify(context.Background(), "user@example.com", strings.Repeat("z", 64)); !errors.Is(err, domain.ErrInvalidVerificationCode) {
		t.Fatalf("expected invalid code, got %v", err)
	}
}

func TestVerifyPropagatesActivationErrors(t *testing.T) {
	expected := errors.New("database unavailable")
	service, err := NewVerificationService(&activationRepository{err: expected})
	if err != nil {
		t.Fatal(err)
	}
	rawCode := []byte("12345678901234567890123456789012")
	_, err = service.Verify(context.Background(), "user@example.com", hex.EncodeToString(rawCode))
	if !errors.Is(err, expected) {
		t.Fatalf("expected wrapped activation error, got %v", err)
	}
}

func TestVerifyPropagatesInvalidVerificationErrors(t *testing.T) {
	service, err := NewVerificationService(&activationRepository{err: domain.ErrInvalidVerificationCode})
	if err != nil {
		t.Fatal(err)
	}
	rawCode := []byte("12345678901234567890123456789012")
	_, err = service.Verify(context.Background(), "user@example.com", hex.EncodeToString(rawCode))
	if !errors.Is(err, domain.ErrInvalidVerificationCode) {
		t.Fatalf("expected invalid verification error, got %v", err)
	}
}

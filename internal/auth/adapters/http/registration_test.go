package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/auth/application"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type handlerRepository struct {
	account domain.Account
	err     error
}

func (r handlerRepository) EmailExists(context.Context, string) (bool, error) {
	return false, nil
}
func (r handlerRepository) FindApprovedCompanyBySecret(context.Context, string) (string, error) {
	return "123e4567-e89b-12d3-a456-426614174000", nil
}
func (r handlerRepository) CreateCandidate(context.Context, string, string, []byte, []byte, time.Time) (domain.Account, error) {
	return r.account, r.err
}
func (r handlerRepository) CreateInterviewer(context.Context, string, string, []byte, []byte, string, time.Time) (domain.Account, error) {
	return r.account, r.err
}
func (r handlerRepository) SaveVerificationChallenge(context.Context, ports.VerificationChallenge) error {
	return nil
}

type handlerChecker struct{}

func (handlerChecker) IsCompromised(context.Context, string) (bool, error) { return false, nil }

type handlerSender struct{}

func (handlerSender) SendVerificationCode(context.Context, string, string) error { return nil }

const registrationPath = "/api/v1/auth/registrations"

func newTestHandler(t *testing.T) *RegistrationHandler {
	t.Helper()
	service, err := application.NewRegistrationService(handlerRepository{
		account: domain.Account{
			ID: "candidate-id", Type: domain.AccountTypeCandidate, Email: "user@example.com",
			Status: domain.PendingVerification, CreatedAt: "2026-10-30T20:00:00Z",
		},
	}, handlerChecker{}, handlerSender{})
	if err != nil {
		t.Fatal(err)
	}
	handler, err := NewRegistrationHandler(service)
	if err != nil {
		t.Fatal(err)
	}
	return handler
}

func TestNewRegistrationHandlerRequiresService(t *testing.T) {
	if _, err := NewRegistrationHandler(nil); err == nil {
		t.Fatal("expected missing service error")
	}
}

func TestRegistrationHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		body := `{"accountType":"candidate","email":"user@example.com","password":"safe password","firstName":"Ada","lastName":"Lovelace"}`
		request := httptest.NewRequest(http.MethodPost, registrationPath, strings.NewReader(body))
		recorder := httptest.NewRecorder()
		newTestHandler(t).ServeHTTP(recorder, request)
		if recorder.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", recorder.Code, recorder.Body.String())
		}
		var response registrationResponse
		if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}
		if response.ID != "candidate-id" || response.Status != domain.PendingVerification {
			t.Fatalf("unexpected response: %+v", response)
		}
	})

	t.Run("method not allowed", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, registrationPath, nil)
		recorder := httptest.NewRecorder()
		newTestHandler(t).ServeHTTP(recorder, request)
		assertProblem(t, recorder, http.StatusMethodNotAllowed, "method_not_allowed")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, registrationPath, strings.NewReader("{"))
		recorder := httptest.NewRecorder()
		newTestHandler(t).ServeHTTP(recorder, request)
		assertProblem(t, recorder, http.StatusBadRequest, "invalid_json")
	})

	t.Run("unknown field", func(t *testing.T) {
		body := `{"accountType":"candidate","email":"user@example.com","password":"safe password","firstName":"Ada","lastName":"Lovelace","unknown":true}`
		request := httptest.NewRequest(http.MethodPost, registrationPath, strings.NewReader(body))
		recorder := httptest.NewRecorder()
		newTestHandler(t).ServeHTTP(recorder, request)
		assertProblem(t, recorder, http.StatusBadRequest, "invalid_json")
	})

	t.Run("validation error", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, registrationPath, bytes.NewBufferString(`{"accountType":"candidate"}`))
		recorder := httptest.NewRecorder()
		newTestHandler(t).ServeHTTP(recorder, request)
		assertProblem(t, recorder, http.StatusUnprocessableEntity, "invalid_registration")
	})

	t.Run("internal error", func(t *testing.T) {
		service, err := application.NewRegistrationService(handlerRepository{err: errors.New("db down")}, handlerChecker{}, handlerSender{})
		if err != nil {
			t.Fatal(err)
		}
		handler, err := NewRegistrationHandler(service)
		if err != nil {
			t.Fatal(err)
		}
		body := `{"accountType":"candidate","email":"user@example.com","password":"safe password","firstName":"Ada","lastName":"Lovelace"}`
		request := httptest.NewRequest(http.MethodPost, registrationPath, strings.NewReader(body))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		assertProblem(t, recorder, http.StatusInternalServerError, "internal_error")
	})
}

func assertProblem(t *testing.T, recorder *httptest.ResponseRecorder, status int, code string) {
	t.Helper()
	if recorder.Code != status {
		t.Fatalf("expected status %d, got %d", status, recorder.Code)
	}
	var problem problemDetails
	if err := json.NewDecoder(recorder.Body).Decode(&problem); err != nil {
		t.Fatal(err)
	}
	if problem.Code != code {
		t.Fatalf("expected code %q, got %q", code, problem.Code)
	}
}

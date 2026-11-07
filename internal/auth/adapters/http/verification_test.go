package http

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/auth/application"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

type verificationHandlerRepository struct {
	err error
}

func (r verificationHandlerRepository) Activate(context.Context, string, []byte, time.Time) (domain.Account, error) {
	if r.err != nil {
		return domain.Account{}, r.err
	}
	return domain.Account{
		ID: "candidate-id", Type: domain.AccountTypeCandidate, Email: "user@example.com",
		Status: "active", CreatedAt: "2026-11-07T18:00:00Z",
	}, nil
}

func newVerificationHandler(t *testing.T, repository verificationHandlerRepository) *VerificationHandler {
	t.Helper()
	service, err := application.NewVerificationService(repository)
	if err != nil {
		t.Fatal(err)
	}
	handler, err := NewVerificationHandler(service)
	if err != nil {
		t.Fatal(err)
	}
	return handler
}

func TestNewVerificationHandlerRequiresService(t *testing.T) {
	if _, err := NewVerificationHandler(nil); err == nil {
		t.Fatal("expected missing service error")
	}
}

func TestVerificationHandler(t *testing.T) {
	rawCode := []byte("12345678901234567890123456789012")
	body := `{"email":"user@example.com","code":"` + hex.EncodeToString(rawCode) + `"}`

	t.Run("success", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		newVerificationHandler(t, verificationHandlerRepository{}).ServeHTTP(
			recorder,
			httptest.NewRequest(http.MethodPost, "/v1/auth/verifications", strings.NewReader(body)),
		)
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
		}
		var response registrationResponse
		if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}
		if response.Status != "active" || response.ID != "candidate-id" {
			t.Fatalf("unexpected response: %+v", response)
		}
	})

	t.Run("method not allowed", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		newVerificationHandler(t, verificationHandlerRepository{}).ServeHTTP(
			recorder, httptest.NewRequest(http.MethodGet, "/v1/auth/verifications", nil),
		)
		assertProblem(t, recorder, http.StatusMethodNotAllowed, "method_not_allowed")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		newVerificationHandler(t, verificationHandlerRepository{}).ServeHTTP(
			recorder, httptest.NewRequest(http.MethodPost, "/v1/auth/verifications", strings.NewReader("{")),
		)
		assertProblem(t, recorder, http.StatusBadRequest, "invalid_json")
	})

	t.Run("unknown field", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		newVerificationHandler(t, verificationHandlerRepository{}).ServeHTTP(
			recorder,
			httptest.NewRequest(http.MethodPost, "/v1/auth/verifications", strings.NewReader(body[:len(body)-1]+`,"extra":true}`)),
		)
		assertProblem(t, recorder, http.StatusBadRequest, "invalid_json")
	})

	t.Run("invalid code", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		newVerificationHandler(t, verificationHandlerRepository{err: domain.ErrInvalidVerificationCode}).ServeHTTP(
			recorder, httptest.NewRequest(http.MethodPost, "/v1/auth/verifications", strings.NewReader(body)),
		)
		assertProblem(t, recorder, http.StatusUnprocessableEntity, "invalid_verification_code")
	})

	t.Run("internal error", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		newVerificationHandler(t, verificationHandlerRepository{err: errors.New("database unavailable")}).ServeHTTP(
			recorder, httptest.NewRequest(http.MethodPost, "/v1/auth/verifications", bytes.NewBufferString(body)),
		)
		assertProblem(t, recorder, http.StatusInternalServerError, "internal_error")
	})
}

func TestVerificationHandlerRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	newVerificationHandler(t, verificationHandlerRepository{}).RegisterRoutes(router.Group("/v1"))
	request := httptest.NewRequest(http.MethodPost, "/v1/auth/verifications", strings.NewReader(`{"email":"","code":""}`))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected registered route, got %d", recorder.Code)
	}
}

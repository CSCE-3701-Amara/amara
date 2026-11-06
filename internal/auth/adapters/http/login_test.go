package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"

	"github.com/CSCE-3701-Amara/amara/internal/auth/application"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type handlerCredentials struct {
	record ports.Credentials
	err    error
}

func (f handlerCredentials) FindCredentials(context.Context, string) (ports.Credentials, error) {
	return f.record, f.err
}

type handlerSessions struct {
	err error
}

func (f *handlerSessions) CreateSession(context.Context, ports.Session) error {
	return f.err
}

type handlerMFA struct {
	err error
}

func (f handlerMFA) GetMFASecret(context.Context, string, domain.AccountType) (ports.MFARecord, error) {
	return ports.MFARecord{}, f.err
}
func (f handlerMFA) SaveMFASecret(context.Context, ports.MFARecord) error { return nil }
func (f handlerMFA) EnableMFA(context.Context, string, domain.AccountType) error {
	return nil
}

func newTestLoginHandler(t *testing.T, credentials ports.Credentials, mfaErr, sessionErr error) *LoginHandler {
	t.Helper()
	mfa, err := application.NewMFAService(handlerMFA{err: mfaErr}, [32]byte{1})
	if err != nil {
		t.Fatal(err)
	}
	service, err := application.NewLoginService(
		handlerCredentials{record: credentials},
		&handlerSessions{err: sessionErr},
		mfa,
	)
	if err != nil {
		t.Fatal(err)
	}
	handler, err := NewLoginHandler(service)
	if err != nil {
		t.Fatal(err)
	}
	return handler
}

func testCredentials(accountType domain.AccountType) ports.Credentials {
	salt := []byte("01234567890123456789")
	return ports.Credentials{
		ID: "123e4567-e89b-12d3-a456-426614174000", Email: "user@example.com",
		AccountType: accountType, PasswordHash: argon2.IDKey([]byte("password123"), salt, 3, 64*1024, 1, 32),
		Salt: salt, Status: "active", CompanyApproved: true,
	}
}

func TestNewLoginHandlerRequiresService(t *testing.T) {
	if _, err := NewLoginHandler(nil); err == nil {
		t.Fatal("expected missing service error")
	}
}

func TestLoginHandlerSuccess(t *testing.T) {
	handler := newTestLoginHandler(t, testCredentials(domain.AccountTypeCandidate), ports.ErrMFANotConfigured, nil)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sessions",
		strings.NewReader(`{"email":"user@example.com","password":"password123"}`))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, http.StatusOK, recorder.Body)
	}
	cookie := recorder.Result().Cookies()
	if len(cookie) != 1 || cookie[0].Name != sessionCookieName || !cookie[0].HttpOnly ||
		!cookie[0].Secure || cookie[0].SameSite != http.SameSiteLaxMode {
		t.Fatalf("unexpected session cookie: %+v", cookie)
	}
	var response struct {
		User struct {
			ID    string `json:"id"`
			Email string `json:"email"`
			Role  string `json:"role"`
		} `json:"user"`
		Session struct {
			ExpiresAt string `json:"expiresAt"`
		} `json:"session"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.User.ID == "" || response.User.Email != "user@example.com" ||
		response.User.Role != "candidate" || response.Session.ExpiresAt == "" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestLoginHandlerRejectsRequests(t *testing.T) {
	handler := newTestLoginHandler(t, testCredentials(domain.AccountTypeCandidate), ports.ErrMFANotConfigured, nil)
	tests := []struct {
		name   string
		method string
		body   string
		want   int
	}{
		{"wrong method", http.MethodGet, "", http.StatusMethodNotAllowed},
		{"invalid JSON", http.MethodPost, "{", http.StatusBadRequest},
		{"unknown field", http.MethodPost, `{"email":"user@example.com","password":"password123","extra":true}`, http.StatusBadRequest},
		{"invalid email", http.MethodPost, `{"email":"not-an-email","password":"password123"}`, http.StatusUnprocessableEntity},
		{"missing password", http.MethodPost, `{"email":"user@example.com"}`, http.StatusUnprocessableEntity},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/api/v1/auth/sessions", strings.NewReader(test.body))
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)
			if recorder.Code != test.want {
				t.Fatalf("status = %d, want %d", recorder.Code, test.want)
			}
		})
	}
}

func TestLoginHandlerMapsServiceErrors(t *testing.T) {
	tests := []struct {
		name   string
		record ports.Credentials
		mfaErr error
		want   int
		code   string
	}{
		{"invalid credentials", ports.Credentials{Status: "active"}, ports.ErrMFANotConfigured, http.StatusUnauthorized, "invalid_credentials"},
		{"inactive", func() ports.Credentials {
			record := testCredentials(domain.AccountTypeCandidate)
			record.Status = "pending_verification"
			return record
		}(), ports.ErrMFANotConfigured, http.StatusForbidden, "account_not_active"},
		{"forbidden", func() ports.Credentials {
			record := testCredentials(domain.AccountTypeCandidate)
			record.Banned = true
			return record
		}(), ports.ErrMFANotConfigured, http.StatusForbidden, "account_forbidden"},
		{"MFA required", testCredentials(domain.AccountTypeInterviewer), ports.ErrMFANotConfigured, http.StatusForbidden, "mfa_required"},
		{"invalid MFA", testCredentials(domain.AccountTypeCandidate), domain.ErrInvalidMFACode, http.StatusUnauthorized, "invalid_mfa_code"},
		{"internal", testCredentials(domain.AccountTypeCandidate), errors.New("database unavailable"), http.StatusInternalServerError, "internal_error"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := newTestLoginHandler(t, test.record, test.mfaErr, nil)
			request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sessions",
				strings.NewReader(`{"email":"user@example.com","password":"password123"}`))
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)
			if recorder.Code != test.want {
				t.Fatalf("status = %d, want %d", recorder.Code, test.want)
			}
			var body problemDetails
			if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body.Code != test.code {
				t.Fatalf("code = %q, want %q", body.Code, test.code)
			}
			if test.code == "internal_error" && body.Detail != "An internal server error occurred." {
				t.Errorf("detail = %q, want fixed internal error message", body.Detail)
			}
		})
	}
}

func TestLoginHandlerMapsSessionPersistenceError(t *testing.T) {
	handler := newTestLoginHandler(t, testCredentials(domain.AccountTypeCandidate), ports.ErrMFANotConfigured, errors.New("storage unavailable"))
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sessions",
		strings.NewReader(`{"email":"user@example.com","password":"password123"}`))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	var body problemDetails
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Detail != "An internal server error occurred." {
		t.Fatalf("detail = %q, want fixed internal error message", body.Detail)
	}
}

func TestLoginHandlerRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := newTestLoginHandler(t, testCredentials(domain.AccountTypeCandidate), ports.ErrMFANotConfigured, nil)
	handler.RegisterRoutes(router.Group("/api/v1"))

	for _, path := range []string{"/api/v1/auth/sessions", "/api/v1/auth/login"} {
		request := httptest.NewRequest(http.MethodPost, path,
			strings.NewReader(`{"email":"user@example.com","password":"password123"}`))
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Errorf("%s: status = %d, want %d", path, recorder.Code, http.StatusOK)
		}
	}
}

func TestLoginHelpers(t *testing.T) {
	for _, value := range []string{"user@example.com", " User@Example.com "} {
		if !validLoginEmail(value) {
			t.Errorf("validLoginEmail(%q) = false", value)
		}
	}
	for _, value := range []string{"", "not-an-email", "User <user@example.com>"} {
		if validLoginEmail(value) {
			t.Errorf("validLoginEmail(%q) = true", value)
		}
	}
	if loginRole(domain.AccountTypeCandidate) != "candidate" ||
		loginRole(domain.AccountTypeInterviewer) != "interviewer" ||
		loginRole(domain.AccountTypeCompany) != "company_representative" {
		t.Fatal("unexpected login role mapping")
	}
	if got := parseSessionTime("invalid"); !got.IsZero() {
		t.Fatalf("invalid session time = %v, want zero", got)
	}
	if got := parseSessionTime(time.Unix(1000, 0).UTC().Format(time.RFC3339)); got.IsZero() {
		t.Fatal("expected valid session time")
	}
}

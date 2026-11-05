package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CSCE-3701-Amara/amara/internal/auth/application"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type companyHandlerRepository struct {
	account      domain.Account
	emailExists  bool
	numberExists bool
	err          error
}

func (r companyHandlerRepository) EmailExists(context.Context, string) (bool, error) {
	return r.emailExists, nil
}
func (r companyHandlerRepository) RegistrationNumberExists(context.Context, string) (bool, error) {
	return r.numberExists, nil
}
func (r companyHandlerRepository) CreateCompanyRegistration(context.Context, ports.CompanyRegistrationRecord) (domain.Account, error) {
	return r.account, r.err
}

type companyHandlerChecker struct{ compromised bool }

func (c companyHandlerChecker) IsCompromised(context.Context, string) (bool, error) {
	return c.compromised, nil
}

func newCompanyHandler(t *testing.T, repository companyHandlerRepository) *CompanyRegistrationHandler {
	t.Helper()
	service, err := application.NewCompanyRegistrationService(repository, companyHandlerChecker{})
	if err != nil {
		t.Fatal(err)
	}
	handler, err := NewCompanyRegistrationHandler(service)
	if err != nil {
		t.Fatal(err)
	}
	return handler
}

func TestNewCompanyRegistrationHandlerRequiresService(t *testing.T) {
	if _, err := NewCompanyRegistrationHandler(nil); err == nil {
		t.Fatal("expected missing service error")
	}
}

func TestCompanyRegistrationHandler(t *testing.T) {
	body := `{"companyName":"Amara Ltd","registrationNumber":"REG-123","email":"owner@example.com","password":"safe password"}`
	t.Run("success", func(t *testing.T) {
		handler := newCompanyHandler(t, companyHandlerRepository{account: domain.Account{
			ID: "request-id", Email: "owner@example.com", Status: domain.PendingReview, CreatedAt: "2026-10-30T20:00:00Z",
		}})
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/api/v1/auth/company-registrations", strings.NewReader(body)))
		if recorder.Code != http.StatusCreated || !strings.Contains(recorder.Body.String(), `"accountType":"company"`) {
			t.Fatalf("unexpected response: %d %s", recorder.Code, recorder.Body.String())
		}
	})
	for _, test := range []struct {
		name       string
		method     string
		body       string
		status     int
		code       string
		repository companyHandlerRepository
	}{
		{"method", http.MethodGet, "", http.StatusMethodNotAllowed, "method_not_allowed", companyHandlerRepository{}},
		{"json", http.MethodPost, "{", http.StatusBadRequest, "invalid_json", companyHandlerRepository{}},
		{"validation", http.MethodPost, `{}`, http.StatusUnprocessableEntity, "invalid_registration", companyHandlerRepository{}},
		{"duplicate email", http.MethodPost, body, http.StatusConflict, "company_registration_exists", companyHandlerRepository{emailExists: true}},
		{"duplicate number", http.MethodPost, body, http.StatusConflict, "company_registration_exists", companyHandlerRepository{numberExists: true}},
		{"internal", http.MethodPost, body, http.StatusInternalServerError, "internal_error", companyHandlerRepository{err: errors.New("db down")}},
	} {
		t.Run(test.name, func(t *testing.T) {
			handler := newCompanyHandler(t, test.repository)
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, httptest.NewRequest(test.method, "/api/v1/auth/company-registrations", strings.NewReader(test.body)))
			assertProblem(t, recorder, test.status, test.code)
		})
	}
}

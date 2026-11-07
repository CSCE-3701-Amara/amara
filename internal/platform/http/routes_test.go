package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	authhttp "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/http"
	"github.com/CSCE-3701-Amara/amara/internal/auth/application"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

// registrationDependencies implements the registration ports without a network
// or database connection so the tests exercise the whole HTTP composition.
type registrationDependencies struct {
	requestID         string
	verified          bool
	panicOnEmailCheck bool
}

func (d *registrationDependencies) EmailExists(ctx context.Context, _ string) (bool, error) {
	if d.panicOnEmailCheck {
		panic("registration repository unavailable")
	}
	d.requestID = RequestID(ctx)
	return false, nil
}

func (d *registrationDependencies) FindApprovedCompanyBySecret(context.Context, string) (string, error) {
	return "company-id", nil
}

func (d *registrationDependencies) CreateCandidate(_ context.Context, _, email string, _, _ []byte, _ time.Time) (domain.Account, error) {
	return domain.Account{ID: "candidate-id", Type: domain.AccountTypeCandidate, Email: email, Status: domain.PendingVerification}, nil
}

func (d *registrationDependencies) CreateInterviewer(_ context.Context, _, email string, _, _ []byte, _ string, _ time.Time) (domain.Account, error) {
	return domain.Account{ID: "interviewer-id", Type: domain.AccountTypeInterviewer, Email: email, Status: domain.PendingVerification}, nil
}

func (d *registrationDependencies) SaveVerificationChallenge(context.Context, ports.VerificationChallenge) error {
	return nil
}

func (d *registrationDependencies) RegistrationNumberExists(context.Context, string) (bool, error) {
	return false, nil
}

func (d *registrationDependencies) CreateCompanyRegistration(_ context.Context, record ports.CompanyRegistrationRecord) (domain.Account, error) {
	return domain.Account{ID: "company-request-id", Type: domain.AccountTypeCompany, Email: record.Email, Status: domain.PendingReview}, nil
}

func (d *registrationDependencies) IsCompromised(context.Context, string) (bool, error) {
	return false, nil
}

func (d *registrationDependencies) SendVerificationCode(context.Context, string, string) error {
	d.verified = true
	return nil
}

func registrationRegistrars(t *testing.T, deps *registrationDependencies) []RouteRegistrar {
	t.Helper()
	registration, err := application.NewRegistrationService(deps, deps, deps)
	if err != nil {
		t.Fatal(err)
	}
	registrationHandler, err := authhttp.NewRegistrationHandler(registration)
	if err != nil {
		t.Fatal(err)
	}
	company, err := application.NewCompanyRegistrationService(deps, deps)
	if err != nil {
		t.Fatal(err)
	}
	companyHandler, err := authhttp.NewCompanyRegistrationHandler(company)
	if err != nil {
		t.Fatal(err)
	}
	return []RouteRegistrar{registrationHandler.RegisterRoutes, companyHandler.RegisterRoutes}
}

func TestRegistrationUsesSharedServer(t *testing.T) {
	const prefix = "/v1"
	for _, tc := range []struct {
		name, path, body, accountType string
		verification                  bool
	}{
		{"candidate", "/auth/registrations", `{"accountType":"candidate","email":"user@example.com","password":"safe password","firstName":"Ada","lastName":"Lovelace"}`, "candidate", true},
		{"interviewer", "/auth/registrations", `{"accountType":"interviewer","email":"user@example.com","password":"safe password","firstName":"Ada","lastName":"Lovelace","companySecret":"company-secret"}`, "interviewer", true},
		{"company", "/auth/company-registrations", `{"companyName":"Amara Ltd","registrationNumber":"REG-123","email":"owner@example.com","password":"safe password"}`, "company", false},
	} {
		t.Run(prefix+"/"+tc.name, func(t *testing.T) {
			deps := &registrationDependencies{}
			var logs bytes.Buffer
			server := testServerWithDB(t, testConfig(), &logs, stubDB("down", "ping failed"), registrationRegistrars(t, deps)...)
			request := httptest.NewRequest(http.MethodPost, prefix+tc.path, strings.NewReader(tc.body))
			request.Header.Set("Origin", "http://localhost:5173")
			response := httptest.NewRecorder()
			server.server.Handler.ServeHTTP(response, request)
			if response.Code != http.StatusCreated {
				t.Fatalf("status = %d, want 201: %s", response.Code, response.Body)
			}
			var body struct {
				AccountType string `json:"accountType"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body.AccountType != tc.accountType || deps.verified != tc.verification {
				t.Fatalf("unexpected registration result: %+v, verification sent: %t", body, deps.verified)
			}
			if id := response.Header().Get("X-Request-ID"); id == "" || id != deps.requestID {
				t.Fatal("registration service must receive the middleware request ID")
			}
			if response.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
				t.Fatal("registration must use the shared CORS middleware")
			}
			if !strings.Contains(logs.String(), `"route":"`+prefix+tc.path+`"`) {
				t.Fatalf("registration must use the shared request logger: %s", &logs)
			}
			ready := httptest.NewRecorder()
			server.server.Handler.ServeHTTP(ready, httptest.NewRequest(http.MethodGet, "/health/ready", nil))
			if ready.Code != http.StatusServiceUnavailable {
				t.Fatalf("registration server must use database.Service readiness: %d", ready.Code)
			}
		})
	}
}

func TestRegistrationUsesSharedRecovery(t *testing.T) {
	deps := &registrationDependencies{panicOnEmailCheck: true}
	var logs bytes.Buffer
	server := testServer(t, testConfig(), &logs, registrationRegistrars(t, deps)...)
	request := httptest.NewRequest(http.MethodPost, "/v1/auth/registrations", strings.NewReader(
		`{"accountType":"candidate","email":"user@example.com","password":"safe password","firstName":"Ada","lastName":"Lovelace"}`,
	))
	response := httptest.NewRecorder()
	server.server.Handler.ServeHTTP(response, request)
	if response.Code != http.StatusInternalServerError || !strings.Contains(response.Body.String(), `"code":"internal_error"`) {
		t.Fatalf("unexpected recovery response: %d %s", response.Code, response.Body)
	}
	if strings.Contains(response.Body.String(), "registration repository unavailable") {
		t.Fatal("panic details must not reach the client")
	}
}

func TestUnavailableRegistrationUsesSharedServer(t *testing.T) {
	var logs bytes.Buffer
	server := testServer(t, testConfig(), &logs, authhttp.RegisterUnavailableRegistrationRoutes)
	const prefix = "/v1"
	request := httptest.NewRequest(http.MethodPost, prefix+"/auth/registrations", strings.NewReader(
		`{"accountType":"candidate","email":"user@example.com","password":"safe password","firstName":"Ada","lastName":"Lovelace"}`,
	))
	request.Header.Set("Origin", "http://localhost:5173")
	response := httptest.NewRecorder()
	server.server.Handler.ServeHTTP(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503: %s", response.Code, response.Body)
	}
	var problem struct {
		Status int    `json:"status"`
		Code   string `json:"code"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &problem); err != nil {
		t.Fatal(err)
	}
	if problem.Status != http.StatusServiceUnavailable || problem.Code != "registration_unavailable" || response.Header().Get("Content-Type") != "application/problem+json" {
		t.Fatalf("unexpected unavailable response: %s", response.Body)
	}
	if response.Header().Get("X-Request-ID") == "" || response.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Fatal("unavailable registration must still receive the shared middleware")
	}
}

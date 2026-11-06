package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/auth/application"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type accountActionTerminator struct {
	accountID string
	token     string
	action    domain.AccountAction
	err       error
}

func (r *accountActionTerminator) PerformAccountAction(_ context.Context, accountID, token string, action domain.AccountAction) error {
	r.accountID, r.token, r.action = accountID, token, action
	return r.err
}

func newTestAccountActionHandler(t *testing.T, repository ports.AccountActionRepository, validator ports.SessionValidator) *AccountActionHandler {
	t.Helper()
	service, err := application.NewAccountActionService(repository)
	if err != nil {
		t.Fatal(err)
	}
	handler, err := NewAccountActionHandler(service, validator, nil)
	if err != nil {
		t.Fatal(err)
	}
	return handler
}

func TestAccountActionHandlerCompletesAndClearsSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repository := &accountActionTerminator{}
	handler := newTestAccountActionHandler(t, repository, sessionValidatorFunc(func(_ context.Context, token string) (ports.Identity, error) {
		if token != "session-token" {
			t.Fatalf("validator token = %q", token)
		}
		return ports.Identity{AccountID: "candidate-id", AccountType: domain.AccountTypeCandidate}, nil
	}))
	router := gin.New()
	handler.RegisterRoutes(router.Group("/v1"))

	request := httptest.NewRequest(http.MethodPost, "/v1/auth/account-actions", strings.NewReader(`{"action":"deactivate"}`))
	request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "session-token"})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent || repository.accountID != "candidate-id" ||
		repository.token != "session-token" || repository.action != domain.AccountActionDeactivate {
		t.Fatalf("status = %d, repository = %+v", recorder.Code, repository)
	}
	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != sessionCookieName || cookies[0].MaxAge >= 0 {
		t.Fatalf("unexpected cleared cookie: %+v", cookies)
	}
}

func TestAccountActionHandlerValidatesRequestsAndAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repository := &accountActionTerminator{}
	handler := newTestAccountActionHandler(t, repository, sessionValidatorFunc(func(context.Context, string) (ports.Identity, error) {
		return ports.Identity{AccountID: "candidate-id", AccountType: domain.AccountTypeCandidate}, nil
	}))
	router := gin.New()
	handler.RegisterRoutes(router.Group("/v1"))

	tests := []struct {
		name string
		body string
		want int
	}{
		{"invalid json", "{", http.StatusBadRequest},
		{"unknown field", `{"action":"delete","extra":true}`, http.StatusBadRequest},
		{"invalid action", `{"action":"archive"}`, http.StatusUnprocessableEntity},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/v1/auth/account-actions", strings.NewReader(test.body))
			request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "session-token"})
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			if recorder.Code != test.want {
				t.Fatalf("status = %d, want %d", recorder.Code, test.want)
			}
		})
	}

	request := httptest.NewRequest(http.MethodPost, "/v1/auth/account-actions", strings.NewReader(`{"action":"delete"}`))
	request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "session-token"})
	router = gin.New()
	handler = newTestAccountActionHandler(t, &accountActionTerminator{err: errors.New("storage down")}, sessionValidatorFunc(func(context.Context, string) (ports.Identity, error) {
		return ports.Identity{AccountID: "candidate-id", AccountType: domain.AccountTypeCandidate}, nil
	}))
	handler.RegisterRoutes(router.Group("/v1"))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("persistence failure status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
}

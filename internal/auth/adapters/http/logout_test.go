package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/auth/application"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type logoutTerminator struct {
	token string
	err   error
}

func (t *logoutTerminator) TerminateSession(_ context.Context, token string) error {
	t.token = token
	return t.err
}

func newTestLogoutHandler(t *testing.T, terminator ports.SessionTerminator, validator ports.SessionValidator) *LogoutHandler {
	t.Helper()
	service, err := application.NewLogoutService(terminator)
	if err != nil {
		t.Fatal(err)
	}
	handler, err := NewLogoutHandler(service, validator, nil)
	if err != nil {
		t.Fatal(err)
	}
	return handler
}

func TestLogoutHandlerTerminatesSessionAndClearsCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	terminator := &logoutTerminator{}
	router := gin.New()
	validator := sessionValidatorFunc(func(_ context.Context, token string) (ports.Identity, error) {
		if token != "session-token" {
			t.Fatalf("validator token = %q", token)
		}
		return ports.Identity{AccountID: "candidate", AccountType: domain.AccountTypeCandidate}, nil
	})
	handler := newTestLogoutHandler(t, terminator, validator)
	handler.RegisterRoutes(router.Group("/v1"))

	request := httptest.NewRequest(http.MethodDelete, "/v1/auth/sessions/current", nil)
	request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "session-token"})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent || terminator.token != "session-token" {
		t.Fatalf("status = %d, token = %q", recorder.Code, terminator.token)
	}
	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != sessionCookieName || cookies[0].MaxAge >= 0 {
		t.Fatalf("unexpected cleared cookie: %+v", cookies)
	}
}

func TestLogoutHandlerRetriesExpiredSessionAndMapsFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validationCount := 0
	handler := newTestLogoutHandler(t, &logoutTerminator{}, sessionValidatorFunc(func(context.Context, string) (ports.Identity, error) {
		validationCount++
		if validationCount == 1 {
			return ports.Identity{AccountID: "candidate", AccountType: domain.AccountTypeCandidate}, nil
		}
		return ports.Identity{}, domain.ErrInvalidSession
	}))
	router := gin.New()
	handler.RegisterRoutes(router.Group("/v1"))
	request := httptest.NewRequest(http.MethodDelete, "/v1/auth/sessions/current", nil)
	request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "expired"})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("first logout status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
	if len(recorder.Result().Cookies()) != 1 {
		t.Fatal("first logout did not clear the session cookie")
	}

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("retry logout status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
	if len(recorder.Result().Cookies()) != 1 {
		t.Fatal("retry logout did not clear the session cookie")
	}

	handler = newTestLogoutHandler(t, &logoutTerminator{err: errors.New("storage down")}, sessionValidatorFunc(func(context.Context, string) (ports.Identity, error) {
		return ports.Identity{AccountID: "candidate", AccountType: domain.AccountTypeCandidate}, nil
	}))
	request = httptest.NewRequest(http.MethodDelete, "/v1/auth/sessions/current", nil)
	request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "session-token"})
	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("persistence failure status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
}

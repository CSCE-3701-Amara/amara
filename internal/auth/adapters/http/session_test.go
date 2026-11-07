package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type sessionValidatorFunc func(context.Context, string) (ports.Identity, error)

func (f sessionValidatorFunc) ValidateSession(ctx context.Context, token string) (ports.Identity, error) {
	return f(ctx, token)
}

func TestCandidateAuthAcceptsCandidateCookieAndPreservesScope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	want := ports.Identity{
		AccountID: "candidate-id", AccountType: domain.AccountTypeCandidate,
		CompanyID: "", ExpiresAt: time.Now().Add(time.Hour),
	}
	router := gin.New()
	router.GET("/feed", CandidateAuth(sessionValidatorFunc(func(_ context.Context, token string) (ports.Identity, error) {
		if token != "session-token" {
			t.Fatalf("validator token = %q", token)
		}
		return want, nil
	})), func(c *gin.Context) {
		got, ok := Identity(c)
		if !ok || got != want {
			t.Fatalf("identity = %#v, %v; want %#v, true", got, ok, want)
		}
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/feed", nil)
	request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "session-token"})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
}

func TestCandidateAuthRejectsNonCandidateRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/feed", CandidateAuth(sessionValidatorFunc(func(context.Context, string) (ports.Identity, error) {
		return ports.Identity{
			AccountID: "interviewer-id", AccountType: domain.AccountTypeInterviewer,
			CompanyID: "company-id",
		}, nil
	})), func(c *gin.Context) {
		t.Fatal("protected handler should not run")
	})

	request := httptest.NewRequest(http.MethodGet, "/feed", nil)
	request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "session-token"})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
}

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

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
	"github.com/gin-gonic/gin"
)

func TestCurrentSessionReturnsValidatedIdentityWithoutToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	expires := time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC)
	for _, role := range []domain.AccountType{domain.AccountTypeCandidate, domain.AccountTypeCompany, domain.AccountTypeInterviewer} {
		t.Run(string(role), func(t *testing.T) {
			calls := 0
			router := gin.New()
			RegisterCurrentSessionRoutes(router.Group("/v1"), sessionValidatorFunc(func(_ context.Context, token string) (ports.Identity, error) {
				calls++
				if token != "private-session-token" {
					t.Fatalf("unexpected session token")
				}
				return ports.Identity{AccountID: "account-id", AccountType: role, ExpiresAt: expires}, nil
			}))
			request := httptest.NewRequest(http.MethodGet, "/v1/auth/sessions/current", nil)
			request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "private-session-token"})
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != 200 || calls != 1 {
				t.Fatalf("status=%d, validator calls=%d", response.Code, calls)
			}
			var got struct {
				User    struct{ ID, Role string }
				Session struct{ ExpiresAt string }
			}
			if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if got.User.ID != "account-id" || got.User.Role != loginRole(role) || got.Session.ExpiresAt != expires.Format(time.RFC3339) {
				t.Fatalf("unexpected identity: %+v", got)
			}
			if response.Header().Get("Cache-Control") != "no-store" {
				t.Fatal("session identity must not be cached")
			}
			if strings.Contains(response.Body.String(), "private-session-token") || len(response.Result().Cookies()) != 0 {
				t.Fatal("session lookup must not expose, issue, or extend tokens")
			}
		})
	}
}

func TestCurrentSessionRejectsMissingInvalidAndUnavailableSessions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		name, token string
		err         error
		status      int
		calls       int
	}{
		{"missing", "", nil, 401, 0},
		{"revoked or expired", "invalid", domain.ErrInvalidSession, 401, 1},
		{"storage unavailable", "valid", errors.New("database down"), 500, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			router := gin.New()
			RegisterCurrentSessionRoutes(router.Group("/v1"), sessionValidatorFunc(func(context.Context, string) (ports.Identity, error) { calls++; return ports.Identity{}, tc.err }))
			request := httptest.NewRequest(http.MethodGet, "/v1/auth/sessions/current", nil)
			if tc.token != "" {
				request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: tc.token})
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != tc.status || calls != tc.calls {
				t.Fatalf("status=%d calls=%d", response.Code, calls)
			}
			if response.Header().Get("Cache-Control") != "no-store" {
				t.Fatal("failures must not be cached")
			}
			if strings.Contains(response.Body.String(), "database down") {
				t.Fatal("storage details leaked")
			}
		})
	}
}

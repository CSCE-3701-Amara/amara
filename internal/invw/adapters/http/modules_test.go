package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	authdomain "github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	authports "github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type sessionWithRole authdomain.AccountType

func (s sessionWithRole) ValidateSession(context.Context, string) (authports.Identity, error) {
	return authports.Identity{AccountType: authdomain.AccountType(s)}, nil
}

func TestModuleRoutesRequireCompanyRole(t *testing.T) {
	for _, role := range []authdomain.AccountType{authdomain.AccountTypeCandidate, authdomain.AccountTypeInterviewer} {
		t.Run(string(role), func(t *testing.T) {
			router := gin.New()
			NewModuleHandler(nil, sessionWithRole(role)).RegisterRoutes(router.Group("/v1"))
			for _, path := range []string{"/v1/test-modules/library", "/v1/test-modules/library/problems", "/v1/test-modules/id"} {
				req := httptest.NewRequest(http.MethodGet, path, nil)
				req.AddCookie(&http.Cookie{Name: "amara_session", Value: "test-session"})
				res := httptest.NewRecorder()
				router.ServeHTTP(res, req)
				if res.Code != http.StatusForbidden {
					t.Errorf("%s returned %d, want 403", path, res.Code)
				}
			}
		})
	}
}

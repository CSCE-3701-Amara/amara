package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	authhttp "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/http"
	authdomain "github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	authports "github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type principalSession struct{ identity authports.Identity }

func (s principalSession) ValidateSession(context.Context, string) (authports.Identity, error) {
	return s.identity, nil
}

func TestJobPrincipalUsesValidatedCompanyScope(t *testing.T) {
	const accountID = "00000000-0000-4000-8000-000000000001"
	const companyID = "00000000-0000-4000-8000-000000000002"
	router := gin.New()
	router.GET("/jobs", authhttp.AuthenticatedAuth(principalSession{identity: authports.Identity{
		AccountID: accountID, CompanyID: companyID, AccountType: authdomain.AccountTypeCompany,
	}}), principalFromSession, func(c *gin.Context) {
		gotID, _ := c.Get(AuthenticatedAccountIDKey)
		gotType, _ := c.Get(AuthenticatedAccountTypeKey)
		if gotID != companyID || gotType != AccountTypeCompany {
			t.Errorf("principal = (%v, %v), want (%s, %s)", gotID, gotType, companyID, AccountTypeCompany)
		}
		c.Status(http.StatusNoContent)
	})
	request := httptest.NewRequest(http.MethodGet, "/jobs", nil)
	request.AddCookie(&http.Cookie{Name: "amara_session", Value: "valid-session"})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204: %s", response.Code, response.Body.String())
	}
	unauthenticated := httptest.NewRecorder()
	router.ServeHTTP(unauthenticated, httptest.NewRequest(http.MethodGet, "/jobs", nil))
	if unauthenticated.Code != http.StatusUnauthorized {
		t.Fatalf("missing cookie status = %d, want 401", unauthenticated.Code)
	}
}

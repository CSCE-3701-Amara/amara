package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	authdomain "github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	jobhttp "github.com/CSCE-3701-Amara/amara/internal/job/adapters/http"
)

func TestPipelineRoutesCoexistWithPostingRoutesAndRequireCompany(t *testing.T) {
	router := gin.New()
	v1 := router.Group("/v1")
	jobhttp.NewPostingHandler(nil, nil).RegisterRoutes(v1)
	NewPipelineHandler(nil, sessionWithRole(authdomain.AccountTypeCandidate)).RegisterRoutes(v1)
	for _, route := range []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/v1/jobs/id/pipeline"},
		{http.MethodPost, "/v1/jobs/id/pipeline"},
		{http.MethodGet, "/v1/pipeline-versions/id/stages"},
		{http.MethodPut, "/v1/pipeline-versions/id/stages"},
	} {
		req := httptest.NewRequest(route.method, route.path, nil)
		req.AddCookie(&http.Cookie{Name: "amara_session", Value: "test-session"})
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		if res.Code != http.StatusForbidden || res.Header().Get("Content-Type") != "application/problem+json" {
			t.Errorf("%s %s returned %d %q", route.method, route.path, res.Code, res.Header().Get("Content-Type"))
		}
	}
}

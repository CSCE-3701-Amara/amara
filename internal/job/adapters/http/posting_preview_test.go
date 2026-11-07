package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authdomain "github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	authports "github.com/CSCE-3701-Amara/amara/internal/auth/ports"
	"github.com/CSCE-3701-Amara/amara/internal/job/application"
	"github.com/CSCE-3701-Amara/amara/internal/job/domain"
	"github.com/CSCE-3701-Amara/amara/internal/job/ports"
	"github.com/gin-gonic/gin"
)

type previewSession struct{}

func (previewSession) ValidateSession(_ context.Context, token string) (authports.Identity, error) {
	return authports.Identity{AccountID: token, AccountType: authdomain.AccountTypeCompany}, nil
}

type previewRepository struct {
	ports.PostingManagementRepository
	posting domain.Posting
}

func (r previewRepository) Get(_ context.Context, companyID, id string) (domain.Posting, error) {
	if companyID != r.posting.CompanyID || id != r.posting.ID {
		return domain.Posting{}, domain.ErrPostingNotFound
	}
	return r.posting, nil
}

func TestPostingPreview(t *testing.T) {
	const id = "00000000-0000-0000-0000-000000000001"
	const owner = "00000000-0000-0000-0000-000000000002"
	repo := previewRepository{posting: domain.Posting{ID: id, CompanyID: owner, Title: "Saved draft", Status: domain.PostingDraft}}
	router := gin.New()
	RegisterPostingPreviewRoutes(router.Group("/v1"), application.NewPostingManagement(repo), previewSession{})
	for _, tc := range []struct {
		name, method, session, body, expected string
		status                                int
	}{
		{"saved draft", http.MethodGet, owner, "", `"title":"Saved draft"`, 200},
		{"incomplete edits", http.MethodPost, owner, `{"title":" Edited ","salary_min":100}`, `"title":"Edited"`, 200},
		{"cross company", http.MethodPost, id, `{"title":"Edited"}`, "", 404},
		{"unauthenticated", http.MethodGet, "", "", "", 401},
		{"malformed", http.MethodPost, owner, `{"title":`, "", 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/v1/jobs/"+id+"/preview", strings.NewReader(tc.body))
			if tc.session != "" {
				req.AddCookie(&http.Cookie{Name: "amara_session", Value: tc.session})
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			body := response.Body.String()
			if response.Code != tc.status || !strings.Contains(body, tc.expected) {
				t.Fatalf("got %d %s, want %d containing %q", response.Code, body, tc.status, tc.expected)
			}
			if tc.status == 200 && (!strings.Contains(body, `"salary_min":null`) || response.Header().Get("Cache-Control") != "no-store") {
				t.Fatalf("unexpected preview body or headers: %s %v", body, response.Header())
			}
		})
	}
}

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
	if companyID != r.posting.CompanyID || (id != r.posting.ID && id != previewVisibleID) {
		return domain.Posting{}, domain.ErrPostingNotFound
	}
	posting := r.posting
	posting.ID = id
	visible := id == previewVisibleID
	posting.SalaryVisible = &visible
	return posting, nil
}

const previewVisibleID = "00000000-0000-0000-0000-000000000003"

func TestPostingPreview(t *testing.T) {
	const id = "00000000-0000-0000-0000-000000000001"
	const owner = "00000000-0000-0000-0000-000000000002"
	minSalary, maxSalary := int32(100), int32(200)
	repo := previewRepository{posting: domain.Posting{ID: id, CompanyID: owner, Title: "Saved draft", Status: domain.PostingDraft, SalaryMin: &minSalary, SalaryMax: &maxSalary}}
	router := gin.New()
	RegisterPostingPreviewRoutes(router.Group("/v1"), application.NewPostingManagement(repo), previewSession{})
	for _, tc := range []struct {
		name, method, session, postingID, body, expected string
		status                                           int
		salaryVisible                                    bool
	}{
		{"saved hidden", http.MethodGet, owner, id, "", `"title":"Saved draft"`, 200, false},
		{"saved visible", http.MethodGet, owner, previewVisibleID, "", `"title":"Saved draft"`, 200, true},
		{"edited hidden", http.MethodPost, owner, id, `{"title":" Edited ","salary_min":100,"salary_max":200,"salary_visible":false}`, `"title":"Edited"`, 200, false},
		{"edited visible", http.MethodPost, owner, id, `{"salary_min":100,"salary_max":200,"salary_visible":true}`, "", 200, true},
		{"visibility omitted", http.MethodPost, owner, id, `{"salary_min":100,"salary_max":200}`, "", 200, false},
		{"empty object", http.MethodPost, owner, id, `{}`, `"title":""`, 200, false},
		{"null body", http.MethodPost, owner, id, `null`, `"code":"invalid_request"`, 400, false},
		{"cross company", http.MethodPost, id, id, `{"title":"Edited"}`, "", 404, false},
		{"unauthenticated", http.MethodGet, "", id, "", "", 401, false},
		{"malformed", http.MethodPost, owner, id, `{"title":`, "", 400, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/v1/jobs/"+tc.postingID+"/preview", strings.NewReader(tc.body))
			if tc.session != "" {
				req.AddCookie(&http.Cookie{Name: "amara_session", Value: tc.session})
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			body := response.Body.String()
			if response.Code != tc.status || !strings.Contains(body, tc.expected) {
				t.Fatalf("got %d %s, want %d containing %q", response.Code, body, tc.status, tc.expected)
			}
			if tc.status != http.StatusOK {
				return
			}
			min, max := `"salary_min":null`, `"salary_max":null`
			if tc.salaryVisible {
				min, max = `"salary_min":100`, `"salary_max":200`
			}
			if !strings.Contains(body, min) || !strings.Contains(body, max) || response.Header().Get("Cache-Control") != "no-store" {
				t.Fatalf("unexpected preview body or headers: %s %v", body, response.Header())
			}
		})
	}
}

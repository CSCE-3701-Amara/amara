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

	"github.com/CSCE-3701-Amara/amara/internal/job/application"
	"github.com/CSCE-3701-Amara/amara/internal/job/domain"
	"github.com/gin-gonic/gin"
)

const testPostingID = "4ce85a22-78c1-48bb-bbf6-1b2e06f6ad17"
const testCandidateID = "1a1a77b2-32ed-45db-a437-7f66271e78af"

type applicationRepoStub struct {
	result domain.Application
	err    error
	calls  int
}

func (s *applicationRepoStub) Apply(context.Context, string, string) (domain.Application, error) {
	s.calls++
	return s.result, s.err
}

func applicationRouter(repo *applicationRepoStub, auth gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	service := application.NewApplicationAdmission(repo)
	handler := NewApplicationHandler(service, auth)
	handler.RegisterRoutes(router.Group("/v1"))
	return router
}
func candidatePrincipal(c *gin.Context) {
	c.Set(AuthenticatedAccountIDKey, testCandidateID)
	c.Set(AuthenticatedAccountTypeKey, "candidate")
	c.Next()
}
func applyRequest(router *gin.Engine, id string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/v1/jobs/"+id+"/apply", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestApplyReturnsCreatedApplicationForAuthenticatedCandidate(t *testing.T) {
	want := domain.Application{ID: "application-id", ReferenceNumber: 17, PostingID: testPostingID, Status: domain.ApplicationAdmitted, SubmittedAt: time.Date(2026, 11, 6, 12, 0, 0, 0, time.UTC), FailedAmarat: []string{}}
	repo := &applicationRepoStub{result: want}
	router := applicationRouter(repo, candidatePrincipal)
	rec := applyRequest(router, testPostingID)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s, want 201", rec.Code, rec.Body.String())
	}
	var got domain.Application
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.ID != want.ID || got.Status != want.Status || got.ReferenceNumber != want.ReferenceNumber {
		t.Fatalf("response=%+v, want fields from %+v", got, want)
	}
	if repo.calls != 1 {
		t.Fatalf("repository calls=%d, want 1", repo.calls)
	}
}
func TestApplyFailsClosedWithoutCandidateAuth(t *testing.T) {
	repo := &applicationRepoStub{}
	rec := applyRequest(applicationRouter(repo, nil), testPostingID)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d, want 401", rec.Code)
	}
	if repo.calls != 0 {
		t.Fatalf("repository calls=%d, want 0", repo.calls)
	}
}
func TestApplyRequiresCandidateRole(t *testing.T) {
	auth := func(c *gin.Context) {
		c.Set(AuthenticatedAccountIDKey, testCandidateID)
		c.Set(AuthenticatedAccountTypeKey, "company")
		c.Next()
	}
	repo := &applicationRepoStub{}
	rec := applyRequest(applicationRouter(repo, auth), testPostingID)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status=%d, want 403", rec.Code)
	}
	if repo.calls != 0 {
		t.Fatalf("repository calls=%d, want 0", repo.calls)
	}
}
func TestApplyRejectsInvalidPostingID(t *testing.T) {
	repo := &applicationRepoStub{}
	rec := applyRequest(applicationRouter(repo, candidatePrincipal), "bad-id")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d, want 400", rec.Code)
	}
	if repo.calls != 0 {
		t.Fatalf("repository calls=%d, want 0", repo.calls)
	}
}
func TestApplyMapsAdmissionConflicts(t *testing.T) {
	cases := []struct {
		name string
		err  error
		code string
	}{
		{"full cap", domain.ErrApplicationLimitReached, "application_limit_reached"},
		{"duplicate", domain.ErrDuplicateApplication, "duplicate_application"},
		{"closed", domain.ErrApplicationPostingUnavailable, "posting_closed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &applicationRepoStub{err: tc.err}
			rec := applyRequest(applicationRouter(repo, candidatePrincipal), testPostingID)
			if rec.Code != http.StatusConflict {
				t.Fatalf("status=%d body=%s, want 409", rec.Code, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), tc.code) {
				t.Fatalf("body=%s, want code %q", rec.Body.String(), tc.code)
			}
		})
	}
}
func TestApplyMapsMissingPosting(t *testing.T) {
	repo := &applicationRepoStub{err: domain.ErrPostingNotFound}
	rec := applyRequest(applicationRouter(repo, candidatePrincipal), testPostingID)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status=%d, want 404", rec.Code)
	}
}
func TestWriteApplicationErrorHidesUnexpectedDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/error", func(c *gin.Context) { writeApplicationError(c, errors.New("sensitive database detail")) })
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/error", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d, want 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "sensitive database detail") {
		t.Fatalf("unexpected error leaked: %s", rec.Body.String())
	}
}

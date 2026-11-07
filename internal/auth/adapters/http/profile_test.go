package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/auth/application"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type profileHandlerRepository struct {
	candidate       domain.CandidateProfilePatch
	company         domain.CompanyProfilePatch
	interviewerName string
	updateCalls     int
	err             error
}

func (r *profileHandlerRepository) GetCandidateProfile(_ context.Context, _ string) (domain.CandidateProfile, error) {
	return domain.CandidateProfile{ID: "candidate-id", Name: "Mona", Email: "mona@example.com", CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)}, r.err
}

func (r *profileHandlerRepository) GetCompanyProfile(_ context.Context, _ string, _ bool) (domain.CompanyProfile, error) {
	return domain.CompanyProfile{ID: "company-id", LegalName: "Example", Approved: true}, r.err
}
func (r *profileHandlerRepository) ListCompanyInterviewers(_ context.Context, _ string, _, _ int32) ([]domain.InterviewerProfile, int64, error) {
	return []domain.InterviewerProfile{{ID: "interviewer-id", CompanyID: "company-id", Name: "Mona", Email: "mona@example.com", AccountStatus: "active"}}, 1, r.err
}
func (r *profileHandlerRepository) GetInterviewerProfile(_ context.Context, _ string) (domain.InterviewerProfile, error) {
	return domain.InterviewerProfile{ID: "interviewer-id", CompanyID: "company-id", Name: "Mona", Email: "mona@example.com", AccountStatus: "active"}, r.err
}
func (r *profileHandlerRepository) GetPublicCompanyProfile(_ context.Context, _ string) (domain.CompanyProfile, error) {
	return domain.CompanyProfile{ID: "company-id", LegalName: "Example", Approved: true}, r.err
}

func (r *profileHandlerRepository) UpdateCandidateProfile(_ context.Context, _ string, patch domain.CandidateProfilePatch) (domain.CandidateProfile, error) {
	r.updateCalls++
	r.candidate = patch
	return domain.CandidateProfile{ID: "candidate-id", Name: "Mona", Email: "mona@example.com", CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC), BirthDate: patch.BirthDate}, r.err
}
func (r *profileHandlerRepository) UpdateCompanyProfile(_ context.Context, _ string, patch domain.CompanyProfilePatch) (domain.CompanyProfile, error) {
	r.updateCalls++
	r.company = patch
	return domain.CompanyProfile{ID: "company-id", LegalName: "Example", Approved: true}, r.err
}
func (r *profileHandlerRepository) UpdateInterviewerProfile(_ context.Context, _ string, name string) (domain.InterviewerProfile, error) {
	r.updateCalls++
	r.interviewerName = name
	return domain.InterviewerProfile{ID: "interviewer-id", CompanyID: "company-id", Name: name, Email: "mona@example.com", AccountStatus: "active"}, r.err
}

func newProfileTestHandler(t *testing.T, repository *profileHandlerRepository, identity ports.Identity) (*ProfileHandler, *gin.Engine) {
	t.Helper()
	service, err := application.NewProfileService(repository)
	if err != nil {
		t.Fatal(err)
	}
	handler, err := NewProfileHandler(service, sessionValidatorFunc(func(context.Context, string) (ports.Identity, error) { return identity, nil }))
	if err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	handler.RegisterRoutes(router.Group("/v1"))
	return handler, router
}

func requestProfile(t *testing.T, router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "session-token"})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func TestProfileHandlerUpdatesCandidate(t *testing.T) {
	repository := &profileHandlerRepository{}
	_, router := newProfileTestHandler(t, repository, ports.Identity{AccountID: "candidate-id", AccountType: domain.AccountTypeCandidate})
	recorder := requestProfile(t, router, http.MethodPatch, "/v1/candidates/me", `{"name":"Mona","birthDate":"2000-01-02"}`)
	if recorder.Code != http.StatusOK || repository.candidate.Name == nil || *repository.candidate.Name != "Mona" || repository.candidate.BirthDate == nil {
		t.Fatalf("status = %d, patch = %+v, body = %s", recorder.Code, repository.candidate, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"birthDate":"2000-01-02"`) {
		t.Fatalf("response = %s", recorder.Body.String())
	}
}

func TestProfileHandlerGetsProfiles(t *testing.T) {
	_, router := newProfileTestHandler(t, &profileHandlerRepository{}, ports.Identity{AccountID: "candidate-id", AccountType: domain.AccountTypeCandidate})
	if recorder := requestProfile(t, router, http.MethodGet, "/v1/candidates/me", ""); recorder.Code != http.StatusOK ||
		!strings.Contains(recorder.Body.String(), `"name":"Mona"`) {
		t.Fatalf("candidate profile status = %d, body = %s", recorder.Code, recorder.Body.String())
	}

	repository := &profileHandlerRepository{}
	_, router = newProfileTestHandler(t, repository, ports.Identity{AccountID: "company-id", AccountType: domain.AccountTypeCompany})
	if recorder := requestProfile(t, router, http.MethodGet, "/v1/companies/me", ""); recorder.Code != http.StatusOK ||
		!strings.Contains(recorder.Body.String(), `"legalName":"Example"`) {
		t.Fatalf("company profile status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if recorder := requestProfile(t, router, http.MethodGet, "/v1/companies/company-id", ""); recorder.Code != http.StatusOK {
		t.Fatalf("public company profile status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if recorder := requestProfile(t, router, http.MethodGet, "/v1/companies/me/interviewers?limit=10&offset=0", ""); recorder.Code != http.StatusOK ||
		!strings.Contains(recorder.Body.String(), `"total":1`) {
		t.Fatalf("interviewer list status = %d, body = %s", recorder.Code, recorder.Body.String())
	}

	_, router = newProfileTestHandler(t, &profileHandlerRepository{}, ports.Identity{AccountID: "interviewer-id", AccountType: domain.AccountTypeInterviewer})
	if recorder := requestProfile(t, router, http.MethodGet, "/v1/interviewers/me", ""); recorder.Code != http.StatusOK ||
		!strings.Contains(recorder.Body.String(), `"name":"Mona"`) {
		t.Fatalf("interviewer profile status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}

func TestProfileHandlerUpdatesCompanyAndInterviewer(t *testing.T) {
	repository := &profileHandlerRepository{}
	_, router := newProfileTestHandler(t, repository, ports.Identity{AccountID: "company-id", AccountType: domain.AccountTypeCompany})
	recorder := requestProfile(t, router, http.MethodPatch, "/v1/companies/me", `{"displayName":null,"contactEmail":"owner@example.com"}`)
	if recorder.Code != http.StatusOK || !repository.company.ClearDisplayName || repository.company.ContactEmail == nil {
		t.Fatalf("company status = %d, patch = %+v", recorder.Code, repository.company)
	}

	repository = &profileHandlerRepository{}
	_, router = newProfileTestHandler(t, repository, ports.Identity{AccountID: "interviewer-id", AccountType: domain.AccountTypeInterviewer})
	recorder = requestProfile(t, router, http.MethodPatch, "/v1/interviewers/me", `{"name":" Mona "}`)
	if recorder.Code != http.StatusOK || repository.interviewerName != "Mona" {
		t.Fatalf("interviewer status = %d, name = %q", recorder.Code, repository.interviewerName)
	}
}

func TestProfileHandlerUpdatesAllCompanyFields(t *testing.T) {
	repository := &profileHandlerRepository{}
	_, router := newProfileTestHandler(t, repository, ports.Identity{AccountID: "company-id", AccountType: domain.AccountTypeCompany})
	body := `{"legalName":"Example","displayName":"Example","description":"Description","websiteLink":"https://example.com","logoUrl":"https://example.com/logo","contactPhone":"123","contactEmail":"owner@example.com","linkedinProfile":"https://linkedin.com/company/example"}`
	if recorder := requestProfile(t, router, http.MethodPatch, "/v1/companies/me", body); recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, patch = %+v", recorder.Code, repository.company)
	}
	if repository.company.LegalName == nil || repository.company.LinkedinProfile == nil {
		t.Fatalf("company fields were not decoded: %+v", repository.company)
	}
}

func TestProfileHandlerRejectsMalformedAndUnauthorizedRequests(t *testing.T) {
	repository := &profileHandlerRepository{}
	_, router := newProfileTestHandler(t, repository, ports.Identity{AccountID: "candidate-id", AccountType: domain.AccountTypeCandidate})
	tests := []struct {
		path string
		body string
		want int
	}{
		{"/v1/candidates/me", "{", http.StatusBadRequest},
		{"/v1/candidates/me", "null", http.StatusBadRequest},
		{"/v1/candidates/me", `{"unknown":true}`, http.StatusBadRequest},
		{"/v1/candidates/me", `{"birthDate":"tomorrow"}`, http.StatusUnprocessableEntity},
		{"/v1/interviewers/me", `{"name":"Mona"}`, http.StatusForbidden},
	}

	for _, test := range tests {
		recorder := requestProfile(t, router, http.MethodPatch, test.path, test.body)
		if recorder.Code != test.want {
			t.Errorf("%s: status = %d, want %d", test.path, recorder.Code, test.want)
		}
	}

	_, companyRouter := newProfileTestHandler(t, &profileHandlerRepository{}, ports.Identity{AccountID: "company-id", AccountType: domain.AccountTypeCompany})
	if recorder := requestProfile(t, companyRouter, http.MethodPatch, "/v1/companies/me", `{}`); recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("empty company patch status = %d", recorder.Code)
	}
}

func TestProfileHandlerRejectsTrailingContentWithoutUpdating(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		identity ports.Identity
		body     string
	}{
		{
			name:     "candidate",
			path:     "/v1/candidates/me",
			identity: ports.Identity{AccountID: "candidate-id", AccountType: domain.AccountTypeCandidate},
			body:     `{"name":"New name"}{"email":"ignored@example.com"}`,
		},
		{
			name:     "company",
			path:     "/v1/companies/me",
			identity: ports.Identity{AccountID: "company-id", AccountType: domain.AccountTypeCompany},
			body:     `{"displayName":"New name"}{"contactEmail":"ignored@example.com"}`,
		},
		{
			name:     "interviewer",
			path:     "/v1/interviewers/me",
			identity: ports.Identity{AccountID: "interviewer-id", AccountType: domain.AccountTypeInterviewer},
			body:     `{"name":"New name"}{"email":"ignored@example.com"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &profileHandlerRepository{}
			_, router := newProfileTestHandler(t, repository, test.identity)
			if recorder := requestProfile(t, router, http.MethodPatch, test.path, test.body); recorder.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
			}
			if repository.updateCalls != 0 {
				t.Fatalf("update service called %d times", repository.updateCalls)
			}
		})
	}
}

func TestProfileHandlerMapsServiceErrors(t *testing.T) {
	repository := &profileHandlerRepository{err: domain.ErrProfileNotFound}
	_, router := newProfileTestHandler(t, repository, ports.Identity{AccountID: "candidate-id", AccountType: domain.AccountTypeCandidate})
	if recorder := requestProfile(t, router, http.MethodPatch, "/v1/candidates/me", `{"name":"Mona"}`); recorder.Code != http.StatusNotFound {
		t.Fatalf("not found status = %d", recorder.Code)
	}
	repository.err = errors.New("storage down")
	if recorder := requestProfile(t, router, http.MethodPatch, "/v1/candidates/me", `{"name":"Mona"}`); recorder.Code != http.StatusInternalServerError {
		t.Fatalf("storage error status = %d", recorder.Code)
	}

	repository.err = nil
	_, interviewerRouter := newProfileTestHandler(t, repository, ports.Identity{AccountID: "interviewer-id", AccountType: domain.AccountTypeInterviewer})
	if recorder := requestProfile(t, interviewerRouter, http.MethodPatch, "/v1/interviewers/me", `{"name":" "}`); recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("invalid interviewer status = %d", recorder.Code)
	}
}

func TestProfileHandlerConstructorValidation(t *testing.T) {
	service, _ := application.NewProfileService(&profileHandlerRepository{})
	if _, err := NewProfileHandler(nil, sessionValidatorFunc(nil)); err == nil {
		t.Fatal("expected nil service error")
	}
	if _, err := NewProfileHandler(service, nil); err == nil {
		t.Fatal("expected nil validator error")
	}
}

func TestProfileHandlerRejectsInvalidCompanyFields(t *testing.T) {
	_, router := newProfileTestHandler(t, &profileHandlerRepository{}, ports.Identity{AccountID: "company-id", AccountType: domain.AccountTypeCompany})
	tests := []string{
		`{"legalName":null}`,
		`{"contactEmail":"invalid"}`,
		`{"websiteLink":"invalid"}`,
	}
	for _, body := range tests {
		if recorder := requestProfile(t, router, http.MethodPatch, "/v1/companies/me", body); recorder.Code != http.StatusUnprocessableEntity {
			t.Errorf("%s: status = %d", body, recorder.Code)
		}
	}
}

func TestProfileHandlerCoversCompanyAndInterviewerFailures(t *testing.T) {
	repository := &profileHandlerRepository{err: errors.New("storage down")}
	_, companyRouter := newProfileTestHandler(t, repository, ports.Identity{AccountID: "company-id", AccountType: domain.AccountTypeCompany})
	if recorder := requestProfile(t, companyRouter, http.MethodPatch, "/v1/companies/me", `{"displayName":"Example"}`); recorder.Code != http.StatusInternalServerError {
		t.Fatalf("company storage status = %d", recorder.Code)
	}

	_, interviewerRouter := newProfileTestHandler(t, repository, ports.Identity{AccountID: "interviewer-id", AccountType: domain.AccountTypeInterviewer})
	if recorder := requestProfile(t, interviewerRouter, http.MethodPatch, "/v1/interviewers/me", `{"name":"Mona"}`); recorder.Code != http.StatusInternalServerError {
		t.Fatalf("interviewer storage status = %d", recorder.Code)
	}
}

func TestProfileHandlerReportsMissingIdentity(t *testing.T) {
	service, _ := application.NewProfileService(&profileHandlerRepository{})
	handler, _ := NewProfileHandler(service, sessionValidatorFunc(func(context.Context, string) (ports.Identity, error) {
		return ports.Identity{}, nil
	}))
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPatch, "/v1/candidates/me", strings.NewReader(`{"name":"Mona"}`))
	handler.updateCandidate(context)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", recorder.Code)
	}

	response := httptest.NewRecorder()
	context, _ = gin.CreateTestContext(response)
	context.Request = httptest.NewRequest(http.MethodPatch, "/v1/interviewers/me", strings.NewReader(`{"name":"Mona"}`))
	handler.updateInterviewer(context)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("interviewer missing identity status = %d", response.Code)
	}

	response = httptest.NewRecorder()
	context, _ = gin.CreateTestContext(response)
	context.Request = httptest.NewRequest(http.MethodPatch, "/v1/companies/me", strings.NewReader(`{"displayName":"Example"}`))
	handler.updateCompany(context)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("company missing identity status = %d", response.Code)
	}
}

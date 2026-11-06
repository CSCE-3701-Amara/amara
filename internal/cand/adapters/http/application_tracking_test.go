package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	authdomain "github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	authports "github.com/CSCE-3701-Amara/amara/internal/auth/ports"
	"github.com/CSCE-3701-Amara/amara/internal/cand/domain"
)

type trackingSessionValidator struct{}

func (trackingSessionValidator) ValidateSession(_ context.Context, token string) (authports.Identity, error) {
	switch token {
	case "candidate":
		return authports.Identity{AccountID: "signed-in-candidate", AccountType: authdomain.AccountTypeCandidate}, nil
	case "company":
		return authports.Identity{AccountType: authdomain.AccountTypeCompany}, nil
	default:
		return authports.Identity{}, authdomain.ErrInvalidSession
	}
}

type trackingReader struct {
	candidateID   string
	applicationID string
}

func (r *trackingReader) List(_ context.Context, candidateID string) ([]domain.ApplicationSummary, error) {
	r.candidateID = candidateID
	return []domain.ApplicationSummary{}, nil
}

func (r *trackingReader) Get(_ context.Context, candidateID, applicationID string) (domain.ApplicationDetail, error) {
	r.candidateID = candidateID
	r.applicationID = applicationID
	return domain.ApplicationDetail{}, domain.ErrApplicationNotFound
}

func TestApplicationTrackingUsesCandidateSession(t *testing.T) {
	const applicationID = "b7195eb6-2ea5-4b95-935a-50cece8d8450"
	const history = "/api/v1/applications"
	cases := []struct {
		name, path, session            string
		status                         int
		wantCandidate, wantApplication string
	}{
		{name: "anonymous", path: history, status: http.StatusUnauthorized},
		{name: "company", path: history + "/" + applicationID, session: "company", status: http.StatusForbidden},
		{name: "empty history", path: history + "?candidate_id=other", session: "candidate", status: http.StatusOK, wantCandidate: "signed-in-candidate"},
		{name: "invalid application id", path: history + "/invalid", session: "candidate", status: http.StatusBadRequest},
		{name: "application outside candidate history", path: history + "/" + applicationID + "?candidate_id=other", session: "candidate", status: http.StatusNotFound, wantCandidate: "signed-in-candidate", wantApplication: applicationID},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reader := &trackingReader{}
			router := gin.New()
			NewApplicationTrackingHandler(reader, trackingSessionValidator{}).RegisterRoutes(router.Group("/api/v1"))
			request := httptest.NewRequest(http.MethodGet, tc.path, nil)
			if tc.session != "" {
				request.AddCookie(&http.Cookie{Name: "amara_session", Value: tc.session})
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d, body = %s", response.Code, tc.status, response.Body)
			}
			if reader.candidateID != tc.wantCandidate || reader.applicationID != tc.wantApplication {
				t.Fatalf("reader called with candidate %q and application %q", reader.candidateID, reader.applicationID)
			}
			if tc.status == http.StatusOK && response.Body.String() != `{"applications":[]}` {
				t.Fatalf("unexpected empty history response: %s", response.Body)
			}
		})
	}
}

package http

import (
	"context"
	"encoding/json"
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
	filter        domain.ApplicationFilter
}

func (r *trackingReader) List(_ context.Context, candidateID string, filter domain.ApplicationFilter) (domain.ApplicationPage, error) {
	r.candidateID = candidateID
	r.filter = filter
	data := []domain.ApplicationSummary{}
	if filter.Status == "rejected" {
		data = append(data, domain.ApplicationSummary{ID: "b7195eb6-2ea5-4b95-935a-50cece8d8450", ReferenceNumber: 1, Status: "rejected", Closed: true, IdentityHidden: true})
	}
	return domain.ApplicationPage{Data: data, Pagination: domain.ApplicationPagination{
		Limit: filter.Limit, Offset: filter.Offset, Total: int64(len(data)),
	}}, nil
}

func (r *trackingReader) Get(_ context.Context, candidateID, applicationID string) (domain.ApplicationDetail, error) {
	r.candidateID = candidateID
	r.applicationID = applicationID
	if applicationID == "b7195eb6-2ea5-4b95-935a-50cece8d8451" {
		decision := false
		return domain.ApplicationDetail{
			ApplicationSummary: domain.ApplicationSummary{ID: applicationID, Status: "rejected", Closed: true, IdentityHidden: true},
			FinalDecision:      &decision,
		}, nil
	}
	return domain.ApplicationDetail{}, domain.ErrApplicationNotFound
}

func TestApplicationTrackingUsesCandidateSession(t *testing.T) {
	const applicationID = "b7195eb6-2ea5-4b95-935a-50cece8d8450"
	const history = "/v1/candidates/me/applications"
	const detail = "/v1/applications/"
	const ownedID = "b7195eb6-2ea5-4b95-935a-50cece8d8451"
	cases := []struct {
		name, path, session            string
		status                         int
		wantCandidate, wantApplication string
		wantFilter                     domain.ApplicationFilter
	}{
		{name: "anonymous", path: history, status: http.StatusUnauthorized},
		{name: "company", path: detail + applicationID, session: "company", status: http.StatusForbidden},
		{name: "empty history", path: history + "?candidate_id=other", session: "candidate", status: http.StatusOK, wantCandidate: "signed-in-candidate", wantFilter: domain.ApplicationFilter{Limit: 20}},
		{name: "invalid limit", path: history + "?limit=101", session: "candidate", status: http.StatusBadRequest},
		{name: "invalid offset", path: history + "?offset=-1", session: "candidate", status: http.StatusBadRequest},
		{name: "filtered page", path: history + "?limit=2&offset=1&status=rejected", session: "candidate", status: http.StatusOK, wantCandidate: "signed-in-candidate", wantFilter: domain.ApplicationFilter{Limit: 2, Offset: 1, Status: "rejected"}},
		{name: "owned application", path: detail + ownedID, session: "candidate", status: http.StatusOK, wantCandidate: "signed-in-candidate", wantApplication: ownedID},
		{name: "invalid application id", path: detail + "invalid", session: "candidate", status: http.StatusBadRequest},
		{name: "malformed UUID separators", path: detail + "b7195eb6x2ea5x4b95x935ax50cece8d8450", session: "candidate", status: http.StatusBadRequest},
		{name: "application outside candidate history", path: detail + applicationID + "?candidate_id=other", session: "candidate", status: http.StatusNotFound, wantCandidate: "signed-in-candidate", wantApplication: applicationID},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reader := &trackingReader{}
			router := gin.New()
			NewApplicationTrackingHandler(reader, trackingSessionValidator{}).RegisterRoutes(router.Group("/v1"))
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
			if reader.filter != tc.wantFilter {
				t.Fatalf("reader filter = %+v, want %+v", reader.filter, tc.wantFilter)
			}
			switch {
			case tc.status != http.StatusOK:
				var problem struct {
					Status int
					Code   string
				}
				if err := json.Unmarshal(response.Body.Bytes(), &problem); err != nil ||
					response.Header().Get("Content-Type") != "application/problem+json" ||
					problem.Status != tc.status || problem.Code == "" {
					t.Fatalf("unexpected problem response: %s", response.Body)
				}
			case tc.name == "empty history":
				if response.Body.String() != `{"data":[],"pagination":{"limit":20,"offset":0,"total":0}}` {
					t.Fatalf("unexpected empty history response: %s", response.Body)
				}
			case tc.name == "filtered page":
				var body struct {
					Data       []map[string]any
					Pagination domain.ApplicationPagination
				}
				if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil ||
					len(body.Data) != 1 || body.Data[0]["referenceNumber"] != float64(1) ||
					body.Data[0]["closed"] != true || body.Data[0]["identityHidden"] != true ||
					body.Pagination != (domain.ApplicationPagination{Limit: 2, Offset: 1, Total: 1}) {
					t.Fatalf("unexpected history response: %s", response.Body)
				}
			case tc.name == "owned application":
				var body map[string]any
				if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil ||
					body["id"] != ownedID || body["finalDecision"] != false || body["closed"] != true {
					t.Fatalf("unexpected application response: %s", response.Body)
				}
				for _, field := range []string{"finalFeedback", "rejectionFeedback"} {
					if _, exists := body[field]; exists {
						t.Fatalf("unreleased %s present: %s", field, response.Body)
					}
				}
			}
		})
	}
}

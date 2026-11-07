package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	candhttp "github.com/CSCE-3701-Amara/amara/internal/cand/adapters/http"
	"github.com/CSCE-3701-Amara/amara/internal/invw/application"
	"github.com/CSCE-3701-Amara/amara/internal/invw/domain"
)

type progressRouteRepository struct {
	candidateID   string
	applicationID string
}

func (r *progressRouteRepository) GetByCandidateApplication(_ context.Context, candidateID, applicationID string) (domain.CandidatePipelineProgress, error) {
	r.candidateID, r.applicationID = candidateID, applicationID
	return domain.CandidatePipelineProgress{
		ApplicationID: applicationID,
		Pipeline:      domain.PipelineVersion{ID: "pipeline-id", Version: 1},
		Stages:        []domain.PipelineStage{{ID: "stage-id", Name: "Review", Type: "manual"}},
	}, nil
}

func TestCandidatePipelineProgressCoexistsWithApplicationRoutes(t *testing.T) {
	for _, progressFirst := range []bool{false, true} {
		router := gin.New()
		v1 := router.Group("/v1")
		progress := NewCandidatePipelineProgressHandler(nil, nil)
		registerApplications := func() {
			candhttp.NewApplicationTrackingHandler(nil, nil).RegisterRoutes(v1)
			candhttp.NewAssessmentProgressHandler(nil, nil).RegisterRoutes(v1)
			candhttp.NewAssessmentResultHandler(nil, nil).RegisterRoutes(v1)
			NewAssessmentHandler(nil, nil).RegisterRoutes(v1)
		}
		if progressFirst {
			progress.RegisterRoutes(v1)
			registerApplications()
		} else {
			registerApplications()
			progress.RegisterRoutes(v1)
		}
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/v1/applications/123/pipeline-progress", nil))
		if response.Code != http.StatusUnauthorized {
			t.Fatalf("progressFirst=%t: status = %d, want 401", progressFirst, response.Code)
		}
	}
}

func TestCandidatePipelineProgressReadsApplicationIDFromSharedWildcard(t *testing.T) {
	const candidateID = "00000000-0000-4000-8000-000000000001"
	const applicationID = "00000000-0000-4000-8000-000000000002"
	repo := &progressRouteRepository{}
	router := gin.New()
	progress := NewCandidatePipelineProgressHandler(
		application.NewCandidatePipelineProgress(repo),
		func(c *gin.Context) { c.Set(CandidateIDContextKey, candidateID) },
	)
	progress.RegisterRoutes(router.Group("/v1"))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/v1/applications/"+applicationID+"/pipeline-progress", nil))
	if response.Code != http.StatusOK || repo.candidateID != candidateID || repo.applicationID != applicationID {
		t.Fatalf("status = %d, candidateID = %q, applicationID = %q; want 200 and route IDs", response.Code, repo.candidateID, repo.applicationID)
	}
}

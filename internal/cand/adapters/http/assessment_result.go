package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/cand/application"
	"github.com/CSCE-3701-Amara/amara/internal/cand/domain"
	httpserver "github.com/CSCE-3701-Amara/amara/internal/platform/http"
)

type AssessmentResultHandler struct {
	results       *application.AssessmentResults
	candidateAuth gin.HandlerFunc
}

func NewAssessmentResultHandler(results *application.AssessmentResults, candidateAuth gin.HandlerFunc) *AssessmentResultHandler {
	return &AssessmentResultHandler{results: results, candidateAuth: candidateAuth}
}

func (h *AssessmentResultHandler) RegisterRoutes(group *gin.RouterGroup) {
	route := group.Group("/applications/:application_id/assessments")
	if h.candidateAuth == nil {
		route.GET("/results", func(c *gin.Context) {
			writeAssessmentError(c, http.StatusUnauthorized, "authentication_required", "Candidate authentication is required.")
		})
		return
	}
	route.GET("/results", h.candidateAuth, h.list)
}

func (h *AssessmentResultHandler) list(c *gin.Context) {
	candidateID, ok := authenticatedCandidateID(c)
	if !ok {
		writeAssessmentError(c, http.StatusUnauthorized, "authentication_required", "Candidate authentication is required.")
		return
	}
	page, err := h.results.List(c.Request.Context(), candidateID, c.Param("application_id"))
	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvalidAssessmentResultQuery):
			writeAssessmentError(c, http.StatusBadRequest, "invalid_assessment_result_query", err.Error())
		case errors.Is(err, domain.ErrAssessmentNotFound):
			writeAssessmentError(c, http.StatusNotFound, "application_not_found", "The application could not be found.")
		default:
			slog.ErrorContext(c.Request.Context(), "list candidate assessment results failed",
				"request_id", httpserver.RequestID(c.Request.Context()), "error", err)
			writeAssessmentError(c, http.StatusInternalServerError, "assessment_results_unavailable", "Assessment results could not be loaded.")
		}
		return
	}
	c.JSON(http.StatusOK, page)
}

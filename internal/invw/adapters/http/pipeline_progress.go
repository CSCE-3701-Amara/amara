package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	authhttp "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/http"
	"github.com/CSCE-3701-Amara/amara/internal/invw/application"
	"github.com/CSCE-3701-Amara/amara/internal/invw/domain"
	httpserver "github.com/CSCE-3701-Amara/amara/internal/platform/http"
)

// CandidateIDContextKey is set by candidate-auth middleware to the authenticated
// candidate UUID.
const CandidateIDContextKey = "candidate_id"

type CandidatePipelineProgressHandler struct {
	progress      *application.CandidatePipelineProgressService
	candidateAuth gin.HandlerFunc
}

func NewCandidatePipelineProgressHandler(progress *application.CandidatePipelineProgressService, candidateAuth gin.HandlerFunc) *CandidatePipelineProgressHandler {
	return &CandidatePipelineProgressHandler{progress: progress, candidateAuth: candidateAuth}
}

func (h *CandidatePipelineProgressHandler) RegisterRoutes(group *gin.RouterGroup) {
	route := group.Group("/applications/:application_id/pipeline-progress")
	if h.candidateAuth == nil {
		route.GET("", func(c *gin.Context) {
			writePipelineProgressError(c, http.StatusUnauthorized, "authentication_required", "Candidate authentication is required.")
		})
		return
	}
	route.GET("", h.candidateAuth, h.get)
}

func (h *CandidatePipelineProgressHandler) get(c *gin.Context) {
	candidateID, ok := authenticatedCandidateID(c)
	if !ok {
		writePipelineProgressError(c, http.StatusUnauthorized, "authentication_required", "Candidate authentication is required.")
		return
	}
	progress, err := h.progress.Get(c.Request.Context(), candidateID, c.Param("application_id"))
	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvalidPipelineProgressQuery):
			writePipelineProgressError(c, http.StatusBadRequest, "invalid_pipeline_progress_query", err.Error())
		case errors.Is(err, domain.ErrPipelineProgressNotFound), errors.Is(err, application.ErrPipelineProgressNotFound):
			writePipelineProgressError(c, http.StatusNotFound, "pipeline_progress_not_found", "The application or its published pipeline could not be found.")
		default:
			slog.ErrorContext(c.Request.Context(), "get candidate pipeline progress failed",
				"request_id", httpserver.RequestID(c.Request.Context()), "error", err)
			writePipelineProgressError(c, http.StatusInternalServerError, "pipeline_progress_unavailable", "Candidate pipeline progress could not be loaded.")
		}
		return
	}
	c.JSON(http.StatusOK, progress)
}

func authenticatedCandidateID(c *gin.Context) (string, bool) {
	if identity, ok := authhttp.Identity(c); ok {
		return identity.AccountID, identity.AccountType == "candidate" && identity.AccountID != ""
	}
	value, exists := c.Get(CandidateIDContextKey)
	if !exists {
		return "", false
	}
	candidateID, ok := value.(string)
	return candidateID, ok && candidateID != ""
}

func writePipelineProgressError(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"error":      gin.H{"code": code, "message": message},
		"request_id": httpserver.RequestID(c.Request.Context()),
	})
}

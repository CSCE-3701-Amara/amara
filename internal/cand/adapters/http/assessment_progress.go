package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/cand/application"
	"github.com/CSCE-3701-Amara/amara/internal/cand/domain"
	httpserver "github.com/CSCE-3701-Amara/amara/internal/platform/http"
)

// CandidateIDContextKey is the Gin context key candidate auth middleware sets
// to the authenticated candidate's UUID.
const CandidateIDContextKey = "candidate_id"

type AssessmentProgressHandler struct {
	progress      *application.AssessmentProgress
	candidateAuth gin.HandlerFunc
}

func NewAssessmentProgressHandler(progress *application.AssessmentProgress, candidateAuth gin.HandlerFunc) *AssessmentProgressHandler {
	return &AssessmentProgressHandler{progress: progress, candidateAuth: candidateAuth}
}

func (h *AssessmentProgressHandler) RegisterRoutes(group *gin.RouterGroup) {
	route := group.Group("/applications/:application_id/assessments")
	if h.candidateAuth == nil {
		unauthorized := func(c *gin.Context) {
			writeAssessmentError(c, http.StatusUnauthorized, "authentication_required", "Candidate authentication is required.")
		}
		route.GET("", unauthorized)
		route.PUT("/:assessment_id/progress", unauthorized)
		return
	}
	route.GET("", h.candidateAuth, h.list)
	route.PUT("/:assessment_id/progress", h.candidateAuth, h.save)
}

func (h *AssessmentProgressHandler) list(c *gin.Context) {
	candidateID, ok := authenticatedCandidateID(c)
	if !ok {
		writeAssessmentError(c, http.StatusUnauthorized, "authentication_required", "Candidate authentication is required.")
		return
	}
	page, err := h.progress.List(c.Request.Context(), candidateID, c.Param("application_id"))
	if err != nil {
		h.writeFailure(c, err)
		return
	}
	c.JSON(http.StatusOK, page)
}

func (h *AssessmentProgressHandler) save(c *gin.Context) {
	candidateID, ok := authenticatedCandidateID(c)
	if !ok {
		writeAssessmentError(c, http.StatusUnauthorized, "authentication_required", "Candidate authentication is required.")
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
	var request struct {
		SavedProgress json.RawMessage `json:"saved_progress"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		writeAssessmentError(c, http.StatusBadRequest, "invalid_assessment_progress", "saved_progress must be a JSON object no larger than 1 MiB.")
		return
	}
	item, err := h.progress.Save(c.Request.Context(), candidateID, c.Param("application_id"), c.Param("assessment_id"), request.SavedProgress)
	if err != nil {
		h.writeFailure(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func authenticatedCandidateID(c *gin.Context) (string, bool) {
	value, exists := c.Get(CandidateIDContextKey)
	if !exists {
		return "", false
	}
	candidateID, ok := value.(string)
	return candidateID, ok && candidateID != ""
}

func (h *AssessmentProgressHandler) writeFailure(c *gin.Context, err error) {
	switch {
	case errors.Is(err, application.ErrInvalidAssessmentProgressQuery):
		writeAssessmentError(c, http.StatusBadRequest, "invalid_assessment_progress_query", err.Error())
	case errors.Is(err, domain.ErrAssessmentNotFound):
		writeAssessmentError(c, http.StatusNotFound, "assessment_not_found", "The application or assessment could not be found.")
	case errors.Is(err, domain.ErrAssessmentNotResumable):
		writeAssessmentError(c, http.StatusConflict, "assessment_not_resumable", "This assessment does not allow progress to be saved or resumed.")
	default:
		slog.ErrorContext(c.Request.Context(), "candidate assessment progress request failed",
			"request_id", httpserver.RequestID(c.Request.Context()), "error", err)
		writeAssessmentError(c, http.StatusInternalServerError, "assessment_progress_unavailable", "Assessment progress could not be loaded or saved.")
	}
}

func writeAssessmentError(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"error":      gin.H{"code": code, "message": message},
		"request_id": httpserver.RequestID(c.Request.Context()),
	})
}

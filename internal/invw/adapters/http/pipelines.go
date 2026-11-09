package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	authhttp "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/http"
	authports "github.com/CSCE-3701-Amara/amara/internal/auth/ports"
	"github.com/CSCE-3701-Amara/amara/internal/invw/adapters/database"
	"github.com/CSCE-3701-Amara/amara/internal/invw/domain"
)

type PipelineHandler struct {
	repo     *database.PipelineRepository
	sessions authports.SessionValidator
}

func NewPipelineHandler(repo *database.PipelineRepository, sessions authports.SessionValidator) *PipelineHandler {
	return &PipelineHandler{repo: repo, sessions: sessions}
}

func (h *PipelineHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	company := authhttp.CompanyAuth(h.sessions, writeError)
	v1.GET("/pipeline-versions/:versionId/graph", company, h.graph(false))
	v1.PUT("/pipeline-versions/:versionId/graph", company, h.graph(true))
	v1.GET("/jobs/:id/pipeline", company, h.get)
	v1.POST("/jobs/:id/pipeline", company, h.create)
	v1.GET("/pipeline-versions/:versionId/stages", company, h.stages)
	v1.PUT("/pipeline-versions/:versionId/stages", company, h.replaceStages)
	v1.POST("/pipeline-versions/:versionId/publish", company, h.publish)
}

func (h *PipelineHandler) publish(c *gin.Context) {
	id, ok := pipelineID(c, "versionId")
	if !ok {
		return
	}
	identity, _ := authhttp.Identity(c)
	if err := h.repo.Publish(c.Request.Context(), identity.CompanyID, id); err != nil {
		writePipelineResult(c, http.StatusNoContent, nil, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *PipelineHandler) get(c *gin.Context) {
	id, ok := pipelineID(c, "id")
	if !ok {
		return
	}
	identity, _ := authhttp.Identity(c)
	pipeline, err := h.repo.Get(c.Request.Context(), identity.CompanyID, id)
	writePipelineResult(c, http.StatusOK, pipeline, err)
}

func (h *PipelineHandler) create(c *gin.Context) {
	id, ok := pipelineID(c, "id")
	if !ok {
		return
	}
	identity, _ := authhttp.Identity(c)
	pipeline, err := h.repo.Create(c.Request.Context(), identity.CompanyID, id)
	if err == nil {
		c.Header("Location", "/v1/jobs/"+id+"/pipeline")
	}
	writePipelineResult(c, http.StatusCreated, pipeline, err)
}

func (h *PipelineHandler) stages(c *gin.Context) {
	id, ok := pipelineID(c, "versionId")
	if !ok {
		return
	}
	identity, _ := authhttp.Identity(c)
	stages, err := h.repo.Stages(c.Request.Context(), identity.CompanyID, id)
	writePipelineResult(c, http.StatusOK, stages, err)
}

func (h *PipelineHandler) replaceStages(c *gin.Context) {
	id, ok := pipelineID(c, "versionId")
	if !ok {
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 2<<20)
	var stages []domain.Stage
	if err := c.ShouldBindJSON(&stages); err != nil || stages == nil {
		writeError(c, http.StatusBadRequest, "invalid_pipeline", "An ordered JSON array of stages is required.")
		return
	}
	for _, stage := range stages {
		if !normalizeIDs(stage.ModuleIDs) || (stage.Interview != nil && !normalizeIDs(stage.Interview.InterviewerIDs)) {
			writeError(c, http.StatusUnprocessableEntity, "invalid_pipeline", "Module and interviewer IDs must be UUIDs.")
			return
		}
	}
	identity, _ := authhttp.Identity(c)
	stages, err := h.repo.ReplaceStages(c.Request.Context(), identity.CompanyID, id, stages)
	writePipelineResult(c, http.StatusOK, stages, err)
}

func pipelineID(c *gin.Context, param string) (string, bool) {
	var id pgtype.UUID
	if err := id.Scan(c.Param(param)); err != nil || !id.Valid {
		writeError(c, http.StatusBadRequest, "invalid_pipeline_id", "The resource ID must be a UUID.")
		return "", false
	}
	return id.String(), true
}

func normalizeIDs(ids []string) bool {
	for i, raw := range ids {
		var id pgtype.UUID
		if err := id.Scan(raw); err != nil || !id.Valid {
			return false
		}
		ids[i] = id.String()
	}
	return true
}

func writePipelineResult(c *gin.Context, status int, result any, err error) {
	var invalid *domain.InvalidPipelineError
	switch {
	case err == nil:
		c.JSON(status, result)
	case errors.Is(err, domain.ErrPipelineNotFound):
		writeError(c, http.StatusNotFound, "pipeline_not_found", err.Error())
	case errors.Is(err, domain.ErrPipelineExists):
		writeError(c, http.StatusConflict, "pipeline_exists", err.Error())
	case errors.Is(err, domain.ErrPipelineLocked):
		writeError(c, http.StatusConflict, "pipeline_locked", err.Error())
	case errors.As(err, &invalid):
		writeError(c, http.StatusUnprocessableEntity, "invalid_pipeline", invalid.Reason)
	default:
		slog.ErrorContext(c.Request.Context(), "pipeline request failed", "error", err)
		writeError(c, http.StatusInternalServerError, "internal_error", "The pipeline request failed.")
	}
}

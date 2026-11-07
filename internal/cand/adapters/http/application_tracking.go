package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	authhttp "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/http"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
	"github.com/CSCE-3701-Amara/amara/internal/cand/domain"
	httpserver "github.com/CSCE-3701-Amara/amara/internal/platform/http"
)

type ApplicationTrackingReader interface {
	List(ctx context.Context, candidateID string, filter domain.ApplicationFilter) (domain.ApplicationPage, error)
	Get(ctx context.Context, candidateID, applicationID string) (domain.ApplicationDetail, error)
}

type ApplicationTrackingHandler struct {
	reader   ApplicationTrackingReader
	sessions ports.SessionValidator
}

func NewApplicationTrackingHandler(reader ApplicationTrackingReader, sessions ports.SessionValidator) *ApplicationTrackingHandler {
	return &ApplicationTrackingHandler{reader: reader, sessions: sessions}
}

func (h *ApplicationTrackingHandler) RegisterRoutes(group *gin.RouterGroup) {
	candidateAuth := authhttp.CandidateAuthWithErrorWriter(h.sessions, writeTrackingError)
	group.GET("/candidates/me/applications", candidateAuth, h.list)
	group.GET("/applications/:application_id", candidateAuth, h.get)
}

func (h *ApplicationTrackingHandler) list(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 || limit > 100 {
		writeTrackingError(c, http.StatusBadRequest, "invalid_limit", "limit must be an integer between 1 and 100.")
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		writeTrackingError(c, http.StatusBadRequest, "invalid_offset", "offset must be a nonnegative integer.")
		return
	}
	filter := domain.ApplicationFilter{Limit: limit, Offset: offset, Status: c.Query("status")}
	identity, _ := authhttp.Identity(c)
	page, err := h.reader.List(c.Request.Context(), identity.AccountID, filter)
	if err != nil {
		writeTrackingError(c, http.StatusInternalServerError, "application_history_unavailable", "Application history could not be loaded.")
		return
	}
	c.JSON(http.StatusOK, page)
}

func (h *ApplicationTrackingHandler) get(c *gin.Context) {
	id := c.Param("application_id")
	var uuid pgtype.UUID
	if uuid.Scan(id) != nil || !strings.EqualFold(uuid.String(), id) {
		writeTrackingError(c, http.StatusBadRequest, "invalid_application_id", "id must be a UUID.")
		return
	}
	identity, _ := authhttp.Identity(c)
	detail, err := h.reader.Get(c.Request.Context(), identity.AccountID, id)
	switch {
	case errors.Is(err, domain.ErrApplicationNotFound):
		writeTrackingError(c, http.StatusNotFound, "application_not_found", "The application was not found.")
	case err != nil:
		writeTrackingError(c, http.StatusInternalServerError, "application_unavailable", "The application could not be loaded.")
	default:
		c.JSON(http.StatusOK, detail)
	}
}

func writeTrackingError(c *gin.Context, status int, code, message string) {
	httpserver.WriteProblem(c, httpserver.ProblemDetails{Status: status, Code: code, Detail: message})
}

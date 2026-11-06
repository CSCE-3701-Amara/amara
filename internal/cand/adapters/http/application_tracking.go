package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	authhttp "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/http"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
	"github.com/CSCE-3701-Amara/amara/internal/cand/domain"
)

type ApplicationTrackingReader interface {
	List(ctx context.Context, candidateID string) ([]domain.ApplicationSummary, error)
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
	routes := group.Group("/applications", authhttp.CandidateAuth(h.sessions))
	routes.GET("", h.list)
	routes.GET("/:id", h.get)
}

func (h *ApplicationTrackingHandler) list(c *gin.Context) {
	identity, _ := authhttp.Identity(c)
	applications, err := h.reader.List(c.Request.Context(), identity.AccountID)
	if err != nil {
		writeTrackingError(c, http.StatusInternalServerError, "application_history_unavailable", "Application history could not be loaded.")
		return
	}
	c.JSON(http.StatusOK, gin.H{"applications": applications})
}

func (h *ApplicationTrackingHandler) get(c *gin.Context) {
	id := c.Param("id")
	var uuid pgtype.UUID
	if len(id) != 36 || uuid.Scan(id) != nil {
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
	c.JSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
}

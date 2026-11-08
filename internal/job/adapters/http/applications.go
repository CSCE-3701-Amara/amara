package http

import (
	"errors"
	"github.com/CSCE-3701-Amara/amara/internal/job/application"
	"github.com/CSCE-3701-Amara/amara/internal/job/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ApplicationHandler exposes candidate application submission and admission.
type ApplicationHandler struct {
	service       *application.ApplicationAdmission
	candidateAuth gin.HandlerFunc
}

func NewApplicationHandler(service *application.ApplicationAdmission, candidateAuth gin.HandlerFunc) *ApplicationHandler {
	return &ApplicationHandler{service: service, candidateAuth: candidateAuth}
}
func (h *ApplicationHandler) RegisterRoutes(group *gin.RouterGroup) {
	route := group.Group("/jobs")
	if h == nil || h.candidateAuth == nil {
		route.POST("/:id/apply", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "authentication_required", "message": "Candidate authentication is required."}})
		})
		return
	}
	route.Use(h.candidateAuth)
	route.POST("/:id/apply", h.apply)
}
func (h *ApplicationHandler) apply(c *gin.Context) {
	if !validUUID(c.Param("id")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_posting_id", "message": "id must be a UUID."}})
		return
	}
	candidateID, ok := authenticatedAccount(c)
	if !ok {
		return
	}
	if authenticatedAccountType(c) != "candidate" {
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "forbidden", "message": "Only a candidate may apply to a posting."}})
		return
	}
	result, err := h.service.Apply(c.Request.Context(), c.Param("id"), candidateID)
	if err != nil {
		writeApplicationError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}
func writeApplicationError(c *gin.Context, err error) {
	status, code, message := http.StatusInternalServerError, "application_unavailable", "The application could not be submitted."
	switch {
	case errors.Is(err, domain.ErrPostingNotFound):
		status, code, message = http.StatusNotFound, "posting_not_found", "The posting was not found."
	case errors.Is(err, domain.ErrApplicationPostingUnavailable):
		status, code, message = http.StatusConflict, "posting_closed", "The posting is not accepting applications."
	case errors.Is(err, domain.ErrApplicationLimitReached):
		status, code, message = http.StatusConflict, "application_limit_reached", "The posting has reached its application limit."
	case errors.Is(err, domain.ErrDuplicateApplication):
		status, code, message = http.StatusConflict, "duplicate_application", "You have already applied to this posting."
	case errors.Is(err, domain.ErrCandidateUnavailable):
		status, code, message = http.StatusForbidden, "candidate_unavailable", "This candidate account cannot apply."
	}
	c.JSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
}

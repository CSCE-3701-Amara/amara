package http

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/auth/application"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type LogoutHandler struct {
	service   *application.LogoutService
	validator ports.SessionValidator
	logger    *slog.Logger
}

func NewLogoutHandler(service *application.LogoutService, validator ports.SessionValidator, logger *slog.Logger) (*LogoutHandler, error) {
	if service == nil {
		return nil, errors.New("logout service must not be nil")
	}
	if validator == nil {
		return nil, errors.New("session validator must not be nil")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &LogoutHandler{service: service, validator: validator, logger: logger}, nil
}

func (h *LogoutHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.DELETE("/auth/sessions/current", logoutAuth(h.validator), gin.WrapH(h))
}

func logoutAuth(validator ports.SessionValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		if validator == nil {
			authenticationRequired(c)
			return
		}
		cookie, err := c.Request.Cookie(sessionCookieName)
		if err != nil || cookie.Value == "" {
			authenticationRequired(c)
			return
		}
		identity, err := validator.ValidateSession(c.Request.Context(), cookie.Value)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidSession) {
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "authentication_unavailable", "message": "Authentication is temporarily unavailable."},
			})
			return
		}
		c.Set(identityContextKey, identity)
		c.Next()
	}
}

func (h *LogoutHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		writeProblem(writer, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	cookie, err := request.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		writeProblem(writer, http.StatusUnauthorized, "authentication_required", "Authentication is required.")
		return
	}
	if err := h.service.Logout(request.Context(), cookie.Value); err != nil {
		h.logger.ErrorContext(request.Context(), "logout failed", "error", err)
		writeProblem(writer, http.StatusInternalServerError, "internal_error", "An internal server error occurred.")
		return
	}
	http.SetCookie(writer, &http.Cookie{
		Name: sessionCookieName, Value: "", Path: "/", MaxAge: -1,
		HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode,
		Expires: time.Unix(1, 0),
	})
	writer.WriteHeader(http.StatusNoContent)
}

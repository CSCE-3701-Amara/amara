package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/auth/application"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type AccountActionHandler struct {
	service   *application.AccountActionService
	validator ports.SessionValidator
	logger    *slog.Logger
}

func NewAccountActionHandler(service *application.AccountActionService, validator ports.SessionValidator, logger *slog.Logger) (*AccountActionHandler, error) {
	if service == nil {
		return nil, errors.New("account action service must not be nil")
	}
	if validator == nil {
		return nil, errors.New("session validator must not be nil")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &AccountActionHandler{service: service, validator: validator, logger: logger}, nil
}

func (h *AccountActionHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/auth/account-actions", CandidateAuth(h.validator), func(c *gin.Context) {
		identity, ok := Identity(c)
		if !ok {
			writeProblem(c.Writer, http.StatusInternalServerError, "internal_error", "An internal server error occurred.")
			return
		}
		h.serveHTTP(c.Writer, c.Request, identity)
	})
}

type accountActionRequest struct {
	Action string `json:"action"`
}

func (h *AccountActionHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.serveHTTP(writer, request, ports.Identity{})
}

func (h *AccountActionHandler) serveHTTP(writer http.ResponseWriter, request *http.Request, identity ports.Identity) {
	if request.Method != http.MethodPost {
		writeProblem(writer, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	var input accountActionRequest
	decoder := json.NewDecoder(http.MaxBytesReader(writer, request.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		writeProblem(writer, http.StatusBadRequest, "invalid_json", "request body must be valid JSON")
		return
	}
	action := domain.AccountAction(strings.TrimSpace(input.Action))
	if err := action.Validate(); err != nil {
		writeProblem(writer, http.StatusUnprocessableEntity, "invalid_account_action", "action must be either deactivate or delete")
		return
	}
	if identity.AccountID == "" {
		writeProblem(writer, http.StatusInternalServerError, "internal_error", "An internal server error occurred.")
		return
	}
	cookie, err := request.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		writeProblem(writer, http.StatusUnauthorized, "authentication_required", "Authentication is required.")
		return
	}
	if err := h.service.Perform(request.Context(), identity.AccountID, cookie.Value, action); err != nil {
		if errors.Is(err, domain.ErrInvalidSession) {
			writeProblem(writer, http.StatusUnauthorized, "authentication_required", "Authentication is required.")
			return
		}
		h.logger.ErrorContext(request.Context(), "account action failed", "error", err)
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

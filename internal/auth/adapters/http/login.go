package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/auth/application"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

const sessionCookieName = "amara_session"

type LoginHandler struct {
	service *application.LoginService
}

func NewLoginHandler(service *application.LoginService) (*LoginHandler, error) {
	if service == nil {
		return nil, errors.New("login service must not be nil")
	}
	return &LoginHandler{service: service}, nil
}

// RegisterRoutes mounts the session endpoint under the API version group.
func (h *LoginHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/auth/sessions", gin.WrapH(h))
	group.POST("/auth/login", gin.WrapH(h))
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	MFACode  string `json:"mfaCode"`
}

func (h *LoginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeProblem(writer, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	var input loginRequest
	decoder := json.NewDecoder(http.MaxBytesReader(writer, request.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		writeProblem(writer, http.StatusBadRequest, "invalid_json", "request body must be valid JSON")
		return
	}
	if !validLoginEmail(input.Email) || strings.TrimSpace(input.Password) == "" {
		writeProblem(writer, http.StatusUnprocessableEntity, "invalid_login", "email and password are required")
		return
	}
	session, err := h.service.Login(request.Context(), domain.Login{
		Email: input.Email, Password: input.Password, MFACode: input.MFACode,
	})
	if err != nil {
		status, code := http.StatusInternalServerError, "internal_error"
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			status, code = http.StatusUnauthorized, "invalid_credentials"
		case errors.Is(err, domain.ErrLoginRateLimited):
			writer.Header().Set("Retry-After", "60")
			status, code = http.StatusTooManyRequests, "login_rate_limited"
		case errors.Is(err, domain.ErrAccountInactive):
			status, code = http.StatusForbidden, "account_not_active"
		case errors.Is(err, domain.ErrAccountForbidden):
			status, code = http.StatusForbidden, "account_forbidden"
		case errors.Is(err, domain.ErrMFARequired):
			status, code = http.StatusForbidden, "mfa_required"
		case errors.Is(err, domain.ErrInvalidMFACode):
			status, code = http.StatusUnauthorized, "invalid_mfa_code"
		}
		writeProblem(writer, status, code, strings.TrimPrefix(err.Error(), "invalid credentials: "))
		return
	}
	http.SetCookie(writer, &http.Cookie{
		Name: sessionCookieName, Value: session.Token, Path: "/", HttpOnly: true,
		Secure: true, SameSite: http.SameSiteLaxMode,
		Expires: parseSessionTime(session.ExpiresAt),
	})
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(map[string]any{
		"user":    map[string]string{"id": session.AccountID, "email": session.Email, "role": loginRole(session.AccountType)},
		"session": map[string]string{"expiresAt": session.ExpiresAt},
	})
}

func validLoginEmail(value string) bool {
	value = strings.TrimSpace(value)
	address, err := mail.ParseAddress(value)
	return err == nil && address.Address == value && strings.Contains(value, "@")
}

func loginRole(accountType domain.AccountType) string {
	if accountType == domain.AccountTypeCompany {
		return "company_representative"
	}
	return string(accountType)
}

func parseSessionTime(value string) (result time.Time) {
	result, _ = time.Parse(time.RFC3339, value)
	return result
}

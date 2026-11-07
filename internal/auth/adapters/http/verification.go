package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/auth/application"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

type VerificationHandler struct {
	service *application.VerificationService
}

func NewVerificationHandler(service *application.VerificationService) (*VerificationHandler, error) {
	if service == nil {
		return nil, errors.New("verification service must not be nil")
	}
	return &VerificationHandler{service: service}, nil
}

func (h *VerificationHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/auth/verifications", gin.WrapH(h))
}

type verificationRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (h *VerificationHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeProblem(writer, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	var input verificationRequest
	decoder := json.NewDecoder(http.MaxBytesReader(writer, request.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		writeProblem(writer, http.StatusBadRequest, "invalid_json", "request body must be valid JSON")
		return
	}
	account, err := h.service.Verify(request.Context(), input.Email, input.Code)
	if err != nil {
		status, code := http.StatusInternalServerError, "internal_error"
		switch {
		case errors.Is(err, domain.ErrInvalidVerificationCode):
			status, code = http.StatusUnprocessableEntity, "invalid_verification_code"
		}
		detail := strings.TrimPrefix(err.Error(), "invalid verification code: ")
		if status >= http.StatusInternalServerError {
			detail = "An internal server error occurred."
		}
		writeProblem(writer, status, code, detail)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(registrationResponse{
		ID: account.ID, AccountType: account.Type, Email: account.Email,
		Status: account.Status, CreatedAt: account.CreatedAt,
	})
}

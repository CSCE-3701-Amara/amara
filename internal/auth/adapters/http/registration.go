package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/CSCE-3701-Amara/amara/internal/auth/application"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

type RegistrationHandler struct {
	service *application.RegistrationService
}

func NewRegistrationHandler(service *application.RegistrationService) (*RegistrationHandler, error) {
	if service == nil {
		return nil, errors.New("registration service must not be nil")
	}
	return &RegistrationHandler{service: service}, nil
}

type registrationRequest struct {
	AccountType   domain.AccountType `json:"accountType"`
	Email         string             `json:"email"`
	Password      string             `json:"password"`
	FirstName     string             `json:"firstName"`
	LastName      string             `json:"lastName"`
	CompanySecret string             `json:"companySecret"`
}

type registrationResponse struct {
	ID          string             `json:"id"`
	AccountType domain.AccountType `json:"accountType"`
	Email       string             `json:"email"`
	Status      string             `json:"status"`
	CreatedAt   string             `json:"createdAt"`
}

type problemDetails struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail,omitempty"`
	Code   string `json:"code,omitempty"`
}

func (h *RegistrationHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeProblem(writer, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	var input registrationRequest
	decoder := json.NewDecoder(http.MaxBytesReader(writer, request.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		writeProblem(writer, http.StatusBadRequest, "invalid_json", "request body must be valid JSON")
		return
	}
	account, err := h.service.Register(request.Context(), domain.Registration{
		AccountType: input.AccountType, Email: input.Email, Password: input.Password,
		FirstName: input.FirstName, LastName: input.LastName, CompanySecret: input.CompanySecret,
	})
	if err != nil {
		status, code := http.StatusInternalServerError, "internal_error"
		switch {
		case errors.Is(err, domain.ErrInvalidRegistration):
			status, code = http.StatusUnprocessableEntity, "invalid_registration"
		case errors.Is(err, domain.ErrEmailAlreadyExists):
			status, code = http.StatusConflict, "email_already_registered"
		case errors.Is(err, domain.ErrCompanyNotFound):
			status, code = http.StatusUnprocessableEntity, "company_association_not_verified"
		}
		writeProblem(writer, status, code, strings.TrimPrefix(err.Error(), "invalid registration: "))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(writer).Encode(registrationResponse{
		ID: account.ID, AccountType: account.Type, Email: account.Email,
		Status: account.Status, CreatedAt: account.CreatedAt,
	})
}

func writeProblem(writer http.ResponseWriter, status int, code, detail string) {
	writer.Header().Set("Content-Type", "application/problem+json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(problemDetails{
		Type: "about:blank", Title: http.StatusText(status), Status: status, Code: code, Detail: detail,
	})
}

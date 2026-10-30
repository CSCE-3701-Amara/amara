package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/CSCE-3701-Amara/amara/internal/auth/application"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

type CompanyRegistrationHandler struct {
	service *application.CompanyRegistrationService
}

func NewCompanyRegistrationHandler(service *application.CompanyRegistrationService) (*CompanyRegistrationHandler, error) {
	if service == nil {
		return nil, errors.New("company registration service must not be nil")
	}
	return &CompanyRegistrationHandler{service: service}, nil
}

type companyRegistrationRequest struct {
	CompanyName        string `json:"companyName"`
	RegistrationNumber string `json:"registrationNumber"`
	Email              string `json:"email"`
	Password           string `json:"password"`
}

func (h *CompanyRegistrationHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeProblem(writer, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	var input companyRegistrationRequest
	decoder := json.NewDecoder(http.MaxBytesReader(writer, request.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		writeProblem(writer, http.StatusBadRequest, "invalid_json", "request body must be valid JSON")
		return
	}
	account, err := h.service.Register(request.Context(), domain.CompanyRegistration{
		CompanyName: input.CompanyName, RegistrationNumber: input.RegistrationNumber,
		Email: input.Email, Password: input.Password,
	})
	if err != nil {
		status, code := http.StatusInternalServerError, "internal_error"
		switch {
		case errors.Is(err, domain.ErrInvalidRegistration):
			status, code = http.StatusUnprocessableEntity, "invalid_registration"
		case errors.Is(err, domain.ErrEmailAlreadyExists), errors.Is(err, domain.ErrCompanyRegistrationExists):
			status, code = http.StatusConflict, "company_registration_exists"
		}
		writeProblem(writer, status, code, strings.TrimPrefix(err.Error(), "invalid registration: "))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(writer).Encode(registrationResponse{
		ID: account.ID, AccountType: domain.AccountTypeCompany, Email: account.Email,
		Status: account.Status, CreatedAt: account.CreatedAt,
	})
}

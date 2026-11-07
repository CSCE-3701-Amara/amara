package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/auth/application"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type ProfileHandler struct {
	service   *application.ProfileService
	validator ports.SessionValidator
}

func NewProfileHandler(service *application.ProfileService, validator ports.SessionValidator) (*ProfileHandler, error) {
	if service == nil {
		return nil, errors.New("profile service must not be nil")
	}
	if validator == nil {
		return nil, errors.New("session validator must not be nil")
	}
	return &ProfileHandler{service: service, validator: validator}, nil
}

func (h *ProfileHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.PATCH("/candidates/me", CandidateAuth(h.validator), h.updateCandidate)
	group.PATCH("/companies/me", CompanyAuth(h.validator, writeSessionError), h.updateCompany)
	group.PATCH("/interviewers/me", AuthenticatedAuth(h.validator), h.updateInterviewer)
}

type candidateProfileRequest struct {
	Name              *string `json:"name"`
	PhoneNumber       *string `json:"phoneNumber"`
	GithubProfile     *string `json:"githubProfile"`
	LinkedinProfile   *string `json:"linkedinProfile"`
	PortfolioLink     *string `json:"portfolioLink"`
	Gender            *string `json:"gender"`
	BirthDate         *string `json:"birthDate"`
	Nationality       *string `json:"nationality"`
	Address           *string `json:"address"`
	Biography         *string `json:"biography"`
	YearsOfExperience *int32  `json:"yearsOfExperience"`
	Available         *bool   `json:"available"`
}

type optionalString struct {
	Value   string
	Present bool
	Null    bool
}

func (o *optionalString) UnmarshalJSON(data []byte) error {
	o.Present = true
	if string(data) == "null" {
		o.Null = true
		return nil
	}
	return json.Unmarshal(data, &o.Value)
}

type companyProfileRequest struct {
	LegalName       optionalString `json:"legalName"`
	DisplayName     optionalString `json:"displayName"`
	Description     optionalString `json:"description"`
	WebsiteLink     optionalString `json:"websiteLink"`
	LogoURL         optionalString `json:"logoUrl"`
	ContactPhone    optionalString `json:"contactPhone"`
	ContactEmail    optionalString `json:"contactEmail"`
	LinkedinProfile optionalString `json:"linkedinProfile"`
}

type interviewerProfileRequest struct {
	Name string `json:"name"`
}

func decodeProfileRequest(w http.ResponseWriter, r *http.Request, input any) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	var raw json.RawMessage
	if err := decoder.Decode(&raw); err != nil || bytes.Equal(raw, []byte("null")) {
		writeProblem(w, http.StatusBadRequest, "invalid_json", "request body must be valid JSON")
		return false
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeProblem(w, http.StatusBadRequest, "invalid_json", "request body must contain exactly one JSON object")
		return false
	}
	payloadDecoder := json.NewDecoder(bytes.NewReader(raw))
	payloadDecoder.DisallowUnknownFields()
	if err := payloadDecoder.Decode(input); err != nil {
		writeProblem(w, http.StatusBadRequest, "invalid_json", "request body must be valid JSON")
		return false
	}
	return true
}

func (h *ProfileHandler) updateCandidate(c *gin.Context) {
	identity, ok := Identity(c)
	if !ok {
		writeProblem(c.Writer, 500, "internal_error", "An internal server error occurred.")
		return
	}
	var input candidateProfileRequest
	if !decodeProfileRequest(c.Writer, c.Request, &input) {
		return
	}
	var birth *time.Time
	if input.BirthDate != nil {
		value, err := time.Parse("2006-01-02", *input.BirthDate)
		if err != nil {
			writeProblem(c.Writer, 422, "invalid_profile", "birthDate must be a valid date")
			return
		}
		birth = &value
	}
	profile, err := h.service.UpdateCandidate(c.Request.Context(), identity.AccountID, domain.CandidateProfilePatch{
		Name: input.Name, PhoneNumber: input.PhoneNumber, GithubProfile: input.GithubProfile,
		LinkedinProfile: input.LinkedinProfile, PortfolioLink: input.PortfolioLink, Gender: input.Gender,
		BirthDate: birth, Nationality: input.Nationality, Address: input.Address, Biography: input.Biography,
		YearsOfExperience: input.YearsOfExperience, Available: input.Available,
	})
	if err != nil {
		writeProfileError(c.Writer, err)
		return
	}
	writeJSON(c.Writer, http.StatusOK, makeCandidateProfileResponse(profile))
}

func (h *ProfileHandler) updateCompany(c *gin.Context) {
	identity, ok := Identity(c)
	if !ok {
		writeProblem(c.Writer, 500, "internal_error", "An internal server error occurred.")
		return
	}
	var input companyProfileRequest
	if !decodeProfileRequest(c.Writer, c.Request, &input) {
		return
	}
	patch := domain.CompanyProfilePatch{}
	set := func(o optionalString, value **string, clear *bool) {
		if !o.Present {
			return
		}
		if o.Null {
			*clear = true
		} else {
			v := o.Value
			*value = &v
		}
	}
	if input.LegalName.Present {
		if input.LegalName.Null {
			writeProblem(c.Writer, 422, "invalid_profile", "legalName cannot be null")
			return
		}
		v := input.LegalName.Value
		patch.LegalName = &v
	}
	set(input.DisplayName, &patch.DisplayName, &patch.ClearDisplayName)
	set(input.Description, &patch.Description, &patch.ClearDescription)
	set(input.WebsiteLink, &patch.WebsiteLink, &patch.ClearWebsiteLink)
	set(input.LogoURL, &patch.LogoURL, &patch.ClearLogoURL)
	set(input.ContactPhone, &patch.ContactPhone, &patch.ClearContactPhone)
	set(input.ContactEmail, &patch.ContactEmail, &patch.ClearContactEmail)
	set(input.LinkedinProfile, &patch.LinkedinProfile, &patch.ClearLinkedinProfile)
	profile, err := h.service.UpdateCompany(c.Request.Context(), identity.AccountID, patch)
	if err != nil {
		writeProfileError(c.Writer, err)
		return
	}
	writeJSON(c.Writer, http.StatusOK, companyProfileResponse{
		ID: profile.ID, LegalName: profile.LegalName, DisplayName: profile.DisplayName,
		Description: profile.Description, WebsiteLink: profile.WebsiteLink, LogoURL: profile.LogoURL,
		ContactPhone: profile.ContactPhone, ContactEmail: profile.ContactEmail,
		LinkedinProfile: profile.LinkedinProfile, Approved: profile.Approved,
	})
}

func (h *ProfileHandler) updateInterviewer(c *gin.Context) {
	identity, ok := Identity(c)
	if !ok {
		writeProblem(c.Writer, 500, "internal_error", "An internal server error occurred.")
		return
	}
	if identity.AccountType != domain.AccountTypeInterviewer {
		writeProblem(c.Writer, http.StatusForbidden, "forbidden", "This account cannot access the requested resource.")
		return
	}
	var input interviewerProfileRequest
	if !decodeProfileRequest(c.Writer, c.Request, &input) {
		return
	}
	profile, err := h.service.UpdateInterviewer(c.Request.Context(), identity.AccountID, input.Name)
	if err != nil {
		writeProfileError(c.Writer, err)
		return
	}
	writeJSON(c.Writer, http.StatusOK, interviewerProfileResponse{
		ID: profile.ID, CompanyID: profile.CompanyID, Name: profile.Name, Email: profile.Email,
		AccountStatus: profile.AccountStatus, CreatedAt: profile.CreatedAt, UpdatedAt: profile.UpdatedAt,
	})
}

func writeProfileError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidProfile):
		writeProblem(w, http.StatusUnprocessableEntity, "invalid_profile", "Profile information is invalid.")
	case errors.Is(err, domain.ErrProfileNotFound):
		writeProblem(w, http.StatusNotFound, "profile_not_found", "Profile was not found.")
	default:
		writeProblem(w, http.StatusInternalServerError, "internal_error", "An internal server error occurred.")
	}
}

type candidateProfileResponse struct {
	ID                string     `json:"id"`
	Name              string     `json:"name"`
	Email             string     `json:"email"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt,omitempty"`
	PhoneNumber       string     `json:"phoneNumber,omitempty"`
	GithubProfile     string     `json:"githubProfile,omitempty"`
	LinkedinProfile   string     `json:"linkedinProfile,omitempty"`
	PortfolioLink     string     `json:"portfolioLink,omitempty"`
	Gender            string     `json:"gender,omitempty"`
	BirthDate         *string    `json:"birthDate,omitempty"`
	Nationality       string     `json:"nationality,omitempty"`
	Address           string     `json:"address,omitempty"`
	Biography         string     `json:"biography,omitempty"`
	YearsOfExperience *int       `json:"yearsOfExperience,omitempty"`
	Available         *bool      `json:"available,omitempty"`
}

func makeCandidateProfileResponse(p domain.CandidateProfile) candidateProfileResponse {
	var birthDate *string
	if p.BirthDate != nil {
		value := p.BirthDate.Format("2006-01-02")
		birthDate = &value
	}
	return candidateProfileResponse{ID: p.ID, Name: p.Name, Email: p.Email, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt, PhoneNumber: p.PhoneNumber, GithubProfile: p.GithubProfile, LinkedinProfile: p.LinkedinProfile, PortfolioLink: p.PortfolioLink, Gender: p.Gender, BirthDate: birthDate, Nationality: p.Nationality, Address: p.Address, Biography: p.Biography, YearsOfExperience: p.YearsOfExperience, Available: p.Available}
}

type companyProfileResponse struct {
	ID              string  `json:"id"`
	LegalName       string  `json:"legalName"`
	DisplayName     *string `json:"displayName"`
	Description     *string `json:"description"`
	WebsiteLink     *string `json:"websiteLink"`
	LogoURL         *string `json:"logoUrl"`
	ContactPhone    *string `json:"contactPhone"`
	ContactEmail    *string `json:"contactEmail"`
	LinkedinProfile *string `json:"linkedinProfile"`
	Approved        bool    `json:"approved"`
}

type interviewerProfileResponse struct {
	ID            string     `json:"id"`
	CompanyID     string     `json:"companyId"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	AccountStatus string     `json:"accountStatus"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

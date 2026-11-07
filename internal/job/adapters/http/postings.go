package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/job/application"
	"github.com/CSCE-3701-Amara/amara/internal/job/domain"
)

const (
	// These Gin context keys are the contract AUTH middleware must populate after
	// validating a session. Identity is never accepted from a request body/header.
	AuthenticatedAccountIDKey   = "authenticated_account_id"
	AuthenticatedAccountTypeKey = "authenticated_account_type"
	AccountTypeCompany          = "company"
	AccountTypeSystemAdmin      = "system_admin"
)

// PostingHandler exposes manager-only posting management and manager/admin deletion.
// The injected middleware must authenticate a session and set the principal keys above.
type PostingHandler struct {
	service *application.PostingManagement
	auth    gin.HandlerFunc
}

func NewPostingHandler(service *application.PostingManagement, auth gin.HandlerFunc) *PostingHandler {
	return &PostingHandler{service: service, auth: auth}
}

type postingRequest struct {
	Title               string               `json:"title"`
	Description         *string              `json:"description"`
	Language            *string              `json:"language"`
	Currency            *string              `json:"currency"`
	SeniorityLevel      string               `json:"seniority_level"`
	EmploymentType      string               `json:"employment_type"`
	Location            *string              `json:"location"`
	SalaryMin           *int32               `json:"salary_min"`
	SalaryMax           *int32               `json:"salary_max"`
	SalaryVisible       *bool                `json:"salary_visible"`
	OpenPositions       *int32               `json:"open_positions"`
	ApplicationOpening  *timeJSON            `json:"application_opening"`
	ApplicationDeadline *timeJSON            `json:"application_deadline"`
	MaxApplications     *int32               `json:"max_applications"`
	FinalistTarget      *int32               `json:"finalist_target"`
	Requirements        []domain.Requirement `json:"requirements"`
	Amarat              []domain.Amara       `json:"amarat"`
}

// timeJSON parses nullable RFC3339 timestamp fields from the posting API body.
type timeJSON struct{ time.Time }

func (t *timeJSON) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

func (r postingRequest) posting() domain.Posting {
	p := domain.Posting{
		Title: r.Title, Description: r.Description, Language: r.Language, Currency: r.Currency,
		SeniorityLevel: r.SeniorityLevel, EmploymentType: r.EmploymentType, Location: r.Location,
		SalaryMin: r.SalaryMin, SalaryMax: r.SalaryMax, SalaryVisible: r.SalaryVisible,
		OpenPositions: r.OpenPositions, MaxApplications: r.MaxApplications, FinalistTarget: r.FinalistTarget,
		Requirements: r.Requirements, Amarat: r.Amarat,
	}
	if r.ApplicationOpening != nil {
		p.ApplicationOpening = &r.ApplicationOpening.Time
	}
	if r.ApplicationDeadline != nil {
		p.ApplicationDeadline = &r.ApplicationDeadline.Time
	}
	return p
}

func (h *PostingHandler) RegisterRoutes(group *gin.RouterGroup) {
	route := group.Group("/jobs")
	if h == nil || h.auth == nil {
		// AUTH has no session middleware in the current checkout. Keep all management
		// routes fail-closed until AUTH supplies its authenticated principal middleware.
		unauthorized := func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "authentication_required", "message": "Company manager authentication is required."}})
		}
		route.POST("", unauthorized)
		route.GET("", unauthorized)
		route.GET("/:id", unauthorized)
		route.PUT("/:id", unauthorized)
		route.POST("/:id/publish", unauthorized)
		route.POST("/:id/republish", unauthorized)
		route.POST("/:id/archive", unauthorized)
		route.DELETE("/:id", unauthorized)
		return
	}
	route.Use(h.auth)
	route.POST("", h.create)
	route.GET("", h.list)
	route.GET("/:id", h.get)
	route.PUT("/:id", h.update)
	route.POST("/:id/publish", h.publish)
	route.POST("/:id/republish", h.publish)
	route.POST("/:id/archive", h.archive)
	route.DELETE("/:id", h.delete)
}

func (h *PostingHandler) create(c *gin.Context) {
	companyID, ok := companyManager(c)
	if !ok {
		return
	}
	request, ok := decodePostingRequest(c)
	if !ok {
		return
	}
	posting, err := h.service.Create(c.Request.Context(), companyID, request.posting())
	if err != nil {
		writePostingError(c, err)
		return
	}
	c.JSON(http.StatusCreated, posting)
}

func (h *PostingHandler) list(c *gin.Context) {
	companyID, ok := companyManager(c)
	if !ok {
		return
	}
	postings, err := h.service.List(c.Request.Context(), companyID)
	if err != nil {
		writePostingError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": postings})
}

func (h *PostingHandler) get(c *gin.Context) {
	companyID, ok := companyManager(c)
	if !ok {
		return
	}
	if !validUUID(c.Param("id")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_posting_id", "message": "id must be a UUID."}})
		return
	}
	posting, err := h.service.Get(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		writePostingError(c, err)
		return
	}
	c.JSON(http.StatusOK, posting)
}

func (h *PostingHandler) update(c *gin.Context) {
	companyID, ok := companyManager(c)
	if !ok {
		return
	}
	if !validUUID(c.Param("id")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_posting_id", "message": "id must be a UUID."}})
		return
	}
	request, ok := decodePostingRequest(c)
	if !ok {
		return
	}
	posting, err := h.service.Update(c.Request.Context(), companyID, c.Param("id"), request.posting())
	if err != nil {
		writePostingError(c, err)
		return
	}
	c.JSON(http.StatusOK, posting)
}

func (h *PostingHandler) publish(c *gin.Context) {
	companyID, ok := companyManager(c)
	if !ok {
		return
	}
	if !validUUID(c.Param("id")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_posting_id", "message": "id must be a UUID."}})
		return
	}
	posting, err := h.service.Publish(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		writePostingError(c, err)
		return
	}
	c.JSON(http.StatusOK, posting)
}

func (h *PostingHandler) archive(c *gin.Context) {
	companyID, ok := companyManager(c)
	if !ok {
		return
	}
	if !validUUID(c.Param("id")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_posting_id", "message": "id must be a UUID."}})
		return
	}
	posting, err := h.service.Archive(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		writePostingError(c, err)
		return
	}
	c.JSON(http.StatusOK, posting)
}

func (h *PostingHandler) delete(c *gin.Context) {
	id, ok := authenticatedAccount(c)
	if !ok {
		return
	}
	if !validUUID(c.Param("id")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_posting_id", "message": "id must be a UUID."}})
		return
	}
	typeName, _ := c.Get(AuthenticatedAccountTypeKey)
	accountType, _ := typeName.(string)
	systemAdmin := accountType == AccountTypeSystemAdmin
	if !systemAdmin && accountType != AccountTypeCompany {
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "forbidden", "message": "Only a company manager or system admin may delete a posting."}})
		return
	}
	companyID := ""
	if !systemAdmin {
		companyID = id
	}
	if err := h.service.Delete(c.Request.Context(), companyID, c.Param("id"), systemAdmin); err != nil {
		writePostingError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func decodePostingRequest(c *gin.Context) (postingRequest, bool) {
	var request postingRequest
	decoder := json.NewDecoder(http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_request", "message": "Request body must contain valid posting JSON."}})
		return postingRequest{}, false
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_request", "message": "Request body must contain one JSON object."}})
		return postingRequest{}, false
	}
	return request, true
}

func authenticatedAccount(c *gin.Context) (string, bool) {
	value, exists := c.Get(AuthenticatedAccountIDKey)
	accountID, ok := value.(string)
	if !exists || !ok || strings.TrimSpace(accountID) == "" || !validUUID(accountID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "authentication_required", "message": "A valid authenticated account is required."}})
		return "", false
	}
	return accountID, true
}

func companyManager(c *gin.Context) (string, bool) {
	accountID, ok := authenticatedAccount(c)
	if !ok {
		return "", false
	}
	typeName, _ := c.Get(AuthenticatedAccountTypeKey)
	if accountType, _ := typeName.(string); accountType != AccountTypeCompany {
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "forbidden", "message": "Only a company manager may manage postings."}})
		return "", false
	}
	return accountID, true
}

func writePostingError(c *gin.Context, err error) {
	status, code, message := http.StatusInternalServerError, "posting_unavailable", "The job posting request could not be completed."
	switch {
	case errors.Is(err, domain.ErrInvalidPosting):
		status, code, message = http.StatusBadRequest, "invalid_posting", err.Error()
	case errors.Is(err, domain.ErrPostingNotFound):
		status, code, message = http.StatusNotFound, "posting_not_found", "The job posting was not found."
	case errors.Is(err, domain.ErrInvalidTransition):
		status, code, message = http.StatusConflict, "invalid_posting_transition", err.Error()
	}
	c.JSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
}

func validUUID(value string) bool {
	if len(value) != 36 || value[8] != '-' || value[13] != '-' || value[18] != '-' || value[23] != '-' {
		return false
	}
	for index, char := range value {
		if index == 8 || index == 13 || index == 18 || index == 23 {
			continue
		}
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return false
		}
	}
	return true
}

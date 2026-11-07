package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	authhttp "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/http"
	authports "github.com/CSCE-3701-Amara/amara/internal/auth/ports"
	"github.com/CSCE-3701-Amara/amara/internal/job/application"
	"github.com/CSCE-3701-Amara/amara/internal/job/domain"
)

type postingPreview struct {
	ID                  string               `json:"id"`
	CompanyID           string               `json:"company_id"`
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
	ApplicationOpening  *time.Time           `json:"application_opening"`
	ApplicationDeadline *time.Time           `json:"application_deadline"`
	Requirements        []domain.Requirement `json:"requirements"`
}

// RegisterPostingPreviewRoutes exposes read-only previews of saved content and editor input.
func RegisterPostingPreviewRoutes(group *gin.RouterGroup, service *application.PostingManagement, sessions authports.SessionValidator) {
	auth := authhttp.CompanyAuth(sessions, func(c *gin.Context, status int, code, message string) {
		c.AbortWithStatusJSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
	})
	preview := func(c *gin.Context) {
		companyID, ok := companyManager(c)
		if !ok {
			return
		}
		id := c.Param("id")
		if !validUUID(id) {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_posting_id", "message": "id must be a UUID."}})
			return
		}
		saved, err := service.Get(c.Request.Context(), companyID, id)
		if err != nil {
			writePostingError(c, err)
			return
		}
		content := saved
		if c.Request.Method == http.MethodPost {
			request, ok := decodePostingRequest(c)
			if !ok {
				return
			}
			content = request.posting().Normalized()
		}
		if content.SalaryVisible == nil || !*content.SalaryVisible {
			content.SalaryMin, content.SalaryMax = nil, nil
		}
		c.Header("Cache-Control", "no-store")
		c.JSON(http.StatusOK, postingPreview{
			ID: saved.ID, CompanyID: saved.CompanyID, Title: content.Title,
			Description: content.Description, Language: content.Language, Currency: content.Currency,
			SeniorityLevel: content.SeniorityLevel, EmploymentType: content.EmploymentType, Location: content.Location,
			SalaryMin: content.SalaryMin, SalaryMax: content.SalaryMax, SalaryVisible: content.SalaryVisible,
			OpenPositions: content.OpenPositions, ApplicationOpening: content.ApplicationOpening,
			ApplicationDeadline: content.ApplicationDeadline, Requirements: content.Requirements,
		})
	}
	group.GET("/jobs/:id/preview", auth, preview)
	group.POST("/jobs/:id/preview", auth, preview)
}

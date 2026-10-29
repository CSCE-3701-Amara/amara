package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/cand/application"
)

// JobFeedHandler adapts the candidate feed service to HTTP. CandidateAuth must
// be supplied by AUTH; nil deliberately fails closed.
type JobFeedHandler struct {
	feed          *application.JobFeed
	candidateAuth gin.HandlerFunc
}

func NewJobFeedHandler(feed *application.JobFeed, candidateAuth gin.HandlerFunc) *JobFeedHandler {
	return &JobFeedHandler{feed: feed, candidateAuth: candidateAuth}
}

func (h *JobFeedHandler) RegisterRoutes(group *gin.RouterGroup) {
	route := group.Group("/jobs")
	if h.candidateAuth == nil {
		route.GET("/feed", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "authentication_required", "message": "Candidate authentication is required."}})
		})
		return
	}
	route.GET("/feed", h.candidateAuth, h.getFeed)
}

func (h *JobFeedHandler) getFeed(c *gin.Context) {
	filter := application.FeedFilter{
		Title: c.Query("title"), Language: c.Query("language"), Currency: c.Query("currency"),
		SeniorityLevel: c.Query("seniority_level"), EmploymentType: c.Query("employment_type"),
		Location: c.Query("location"), Cursor: c.Query("cursor"),
	}
	if raw := c.Query("limit"); raw != "" {
		limit, err := strconv.Atoi(raw)
		if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_limit", "message": "limit must be an integer."}}); return }
		filter.Limit = limit
	}
	if raw := c.Query("salary_min"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_salary_min", "message": "salary_min must be a 32-bit integer."}}); return }
		n := int32(value); filter.SalaryMin = &n
	}
	if raw := c.Query("salary_max"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_salary_max", "message": "salary_max must be a 32-bit integer."}}); return }
		n := int32(value); filter.SalaryMax = &n
	}
	page, err := h.feed.List(c.Request.Context(), filter)
	if err != nil {
		if strings.HasPrefix(err.Error(), "limit must be between") || strings.HasPrefix(err.Error(), "invalid cursor") {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_feed_query", "message": err.Error()}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "feed_unavailable", "message": "The job feed could not be loaded."}})
		return
	}
	c.JSON(http.StatusOK, page)
}

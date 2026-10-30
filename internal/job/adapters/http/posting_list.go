package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/job/application"
	httpserver "github.com/CSCE-3701-Amara/amara/internal/platform/http"
)

// companyIDParam names the company wildcard. Every /company/... route must
const companyIDParam = "company_id"

// PostingListHandler adapts the company posting list to HTTP.
type PostingListHandler struct {
	list              *application.CompanyPostingList
	hiringManagerAuth gin.HandlerFunc
}

// until AUTH is implemented.
func NewPostingListHandler(list *application.CompanyPostingList, hiringManagerAuth gin.HandlerFunc) *PostingListHandler {
	return &PostingListHandler{list: list, hiringManagerAuth: hiringManagerAuth}
}

// TODO(AUTH): require hiringManagerAuth and check that :company_id matches the
// session's company. Until then the path is trusted, and an unknown company
// returns an empty list rather than an error.
func (h *PostingListHandler) RegisterRoutes(group *gin.RouterGroup) {
	route := group.Group("/company/:" + companyIDParam)
	if h.hiringManagerAuth != nil {
		route.Use(h.hiringManagerAuth)
	}
	route.GET("/postings", h.listPostings)
}

func (h *PostingListHandler) listPostings(c *gin.Context) {
	filter := application.PostingListFilter{
		CompanyID: c.Param(companyIDParam), Title: c.Query("title"),
		Status: c.Query("status"), Cursor: c.Query("cursor"),
	}
	if raw := c.Query("limit"); raw != "" {
		limit, err := strconv.Atoi(raw)
		if err != nil {
			writeError(c, http.StatusBadRequest, "invalid_limit", "limit must be an integer.")
			return
		}
		filter.Limit = limit
	}
	page, err := h.list.List(c.Request.Context(), filter)
	if err != nil {
		if errors.Is(err, application.ErrInvalidPostingListQuery) {
			message := strings.TrimPrefix(err.Error(), application.ErrInvalidPostingListQuery.Error()+": ")
			writeError(c, http.StatusBadRequest, "invalid_posting_list_query", message)
			return
		}
		slog.ErrorContext(c.Request.Context(), "list company postings failed",
			"request_id", httpserver.RequestID(c.Request.Context()), "error", err)
		writeError(c, http.StatusInternalServerError, "posting_list_unavailable", "The posting list could not be loaded.")
		return
	}
	c.JSON(http.StatusOK, page)
}

func writeError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"error":      gin.H{"code": code, "message": message},
		"request_id": httpserver.RequestID(c.Request.Context()),
	})
}

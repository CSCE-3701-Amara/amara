package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	authhttp "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/http"
	authports "github.com/CSCE-3701-Amara/amara/internal/auth/ports"
	"github.com/CSCE-3701-Amara/amara/internal/job/application"
	httpserver "github.com/CSCE-3701-Amara/amara/internal/platform/http"
)

// companyIDParam names the company wildcard.
const companyIDParam = "company_id"

// PostingListHandler adapts the company posting list to HTTP.
type PostingListHandler struct {
	list     *application.CompanyPostingList
	sessions authports.SessionValidator
}

func NewPostingListHandler(list *application.CompanyPostingList, sessions authports.SessionValidator) *PostingListHandler {
	return &PostingListHandler{list: list, sessions: sessions}
}

func (h *PostingListHandler) RegisterRoutes(group *gin.RouterGroup) {
	route := group.Group("/company/:" + companyIDParam)
	route.Use(authhttp.CompanyAuth(h.sessions, writeError))
	route.GET("/postings", h.listPostings)
}

func (h *PostingListHandler) listPostings(c *gin.Context) {
	identity, ok := authhttp.Identity(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "authentication_required", "Authentication is required.")
		return
	}
	companyID := c.Param(companyIDParam)
	if !strings.EqualFold(identity.CompanyID, companyID) {
		writeError(c, http.StatusForbidden, "forbidden", "You cannot access another company's postings.")
		return
	}

	filter := application.PostingListFilter{
		CompanyID: companyID, Title: c.Query("title"),
		Status: c.Query("status"), Cursor: c.Query("cursor"),
	}
	if raw, supplied := c.Request.URL.Query()["limit"]; supplied {
		if len(raw) != 1 {
			writeError(c, http.StatusBadRequest, "invalid_limit", "limit must be supplied once as an integer between 1 and 100.")
			return
		}
		limit, err := strconv.Atoi(raw[0])
		if err != nil {
			writeError(c, http.StatusBadRequest, "invalid_limit", "limit must be an integer.")
			return
		}
		if limit < 1 || limit > 100 {
			writeError(c, http.StatusBadRequest, "invalid_limit", "limit must be between 1 and 100.")
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

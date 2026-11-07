package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	authhttp "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/http"
	authports "github.com/CSCE-3701-Amara/amara/internal/auth/ports"
	"github.com/CSCE-3701-Amara/amara/internal/invw/adapters/database"
	"github.com/CSCE-3701-Amara/amara/internal/invw/domain"
	httpserver "github.com/CSCE-3701-Amara/amara/internal/platform/http"
)

type ModuleHandler struct {
	repo     *database.ModuleRepository
	sessions authports.SessionValidator
}

func NewModuleHandler(repo *database.ModuleRepository, sessions authports.SessionValidator) *ModuleHandler {
	return &ModuleHandler{repo: repo, sessions: sessions}
}

func (h *ModuleHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	group := v1.Group("/modules", authhttp.CompanyAuth(h.sessions, writeError))
	group.GET("", h.list(false))
	group.GET("/library", h.list(true))
	group.GET("/library/problems", h.problems)
	group.POST("", h.save(false))
	group.GET("/:id", h.moduleAction("get"))
	group.PUT("/:id", h.save(true))
	group.POST("/:id/publish", h.moduleAction("publish"))
	group.POST("/:id/copy", h.moduleAction("copy"))
}

func (h *ModuleHandler) list(library bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		identity, _ := authhttp.Identity(c)
		items, err := h.repo.List(c.Request.Context(), identity.CompanyID, library)
		writeResult(c, http.StatusOK, items, err)
	}
}

func (h *ModuleHandler) problems(c *gin.Context) {
	items, err := h.repo.Problems(c.Request.Context())
	writeResult(c, http.StatusOK, items, err)
}

func (h *ModuleHandler) save(update bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 2<<20)
		var definition domain.Definition
		if err := c.ShouldBindJSON(&definition); err != nil {
			writeError(c, http.StatusBadRequest, "invalid_module", "A valid module definition is required.")
			return
		}
		if err := definition.ValidateDraft(); err != nil {
			writeError(c, http.StatusUnprocessableEntity, "invalid_module", err.Error())
			return
		}
		if definition.Library != nil {
			definition.Library.Pool = nil
		}
		identity, _ := authhttp.Identity(c)
		ctx := c.Request.Context()
		if update {
			id, ok := moduleID(c)
			if !ok {
				return
			}
			module, err := h.repo.Update(ctx, identity.CompanyID, id, definition)
			writeResult(c, http.StatusOK, module, err)
			return
		}
		module, err := h.repo.Create(ctx, identity.CompanyID, definition)
		writeResult(c, http.StatusCreated, module, err)
	}
}

func (h *ModuleHandler) moduleAction(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := moduleID(c)
		if !ok {
			return
		}
		identity, _ := authhttp.Identity(c)
		ctx := c.Request.Context()
		switch action {
		case "get":
			module, err := h.repo.Get(ctx, identity.CompanyID, id)
			writeResult(c, http.StatusOK, module, err)
		case "publish":
			module, err := h.repo.Publish(ctx, identity.CompanyID, id)
			writeResult(c, http.StatusOK, module, err)
		case "copy":
			module, err := h.repo.Copy(ctx, identity.CompanyID, id)
			writeResult(c, http.StatusCreated, module, err)
		}
	}
}

func moduleID(c *gin.Context) (string, bool) {
	id := c.Param("id")
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil || !uuid.Valid {
		writeError(c, http.StatusBadRequest, "invalid_module_id", "The module ID must be a UUID.")
		return "", false
	}
	return uuid.String(), true
}

func writeResult(c *gin.Context, status int, result any, err error) {
	var invalid *domain.InvalidModuleError
	switch {
	case err == nil:
		if module, ok := result.(domain.Module); ok && status == http.StatusCreated {
			c.Header("Location", "/v1/modules/"+module.ID)
		}
		c.JSON(status, result)
	case errors.Is(err, domain.ErrModuleNotFound):
		writeError(c, http.StatusNotFound, "module_not_found", "The requested test module was not found.")
	case errors.Is(err, domain.ErrModuleLocked):
		writeError(c, http.StatusConflict, "module_locked", err.Error())
	case errors.As(err, &invalid):
		writeError(c, http.StatusUnprocessableEntity, "invalid_module", invalid.Reason)
	default:
		slog.ErrorContext(c.Request.Context(), "test module request failed", "error", err)
		writeError(c, http.StatusInternalServerError, "internal_error", "The test module request failed.")
	}
}

func writeError(c *gin.Context, status int, code, message string) {
	httpserver.WriteProblem(c, httpserver.ProblemDetails{Status: status, Code: code, Detail: message})
}

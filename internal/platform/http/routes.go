package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/platform/config"
)

// RouteRegistrar mounts a module's HTTP handlers under /v1. Each module applies
// its own authentication and authorization middleware to protected routes.
type RouteRegistrar func(*gin.RouterGroup)

func newRouter(cfg config.HTTPConfig, logger *slog.Logger, registrars ...RouteRegistrar) (*gin.Engine, error) {
	router := gin.New()
	router.HandleMethodNotAllowed = true
	router.RedirectTrailingSlash = false
	if err := router.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		return nil, fmt.Errorf("configure trusted proxies: %w", err)
	}

	router.Use(requestID(), requestLogger(logger), Recovery(logger), ErrorHandler(logger), corsMiddleware(cfg.AllowedOrigins))
	registerRoutes(router, registrars...)
	return router, nil
}

func registerRoutes(router *gin.Engine, registrars ...RouteRegistrar) {
	// Liveness only: this endpoint does not claim that external providers are ready.
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := router.Group("/v1")
	for _, register := range registrars {
		register(v1)
	}
	router.NoRoute(func(c *gin.Context) {
		WriteProblem(c, ProblemDetails{
			Type:   "about:blank",
			Title:  "Not Found",
			Status: http.StatusNotFound,
			Detail: "The requested resource was not found.",
			Code:   "not_found",
		})
	})
	router.NoMethod(func(c *gin.Context) {
		WriteProblem(c, ProblemDetails{
			Type:   "about:blank",
			Title:  "Method Not Allowed",
			Status: http.StatusMethodNotAllowed,
			Detail: "The request method is not allowed.",
			Code:   "method_not_allowed",
		})
	})
}

func writeError(c *gin.Context, status int, code, message string) {
	WriteProblem(c, ProblemDetails{
		Type:   "about:blank",
		Title:  defaultTitle(status),
		Status: status,
		Detail: message,
		Code:   code,
	})
}

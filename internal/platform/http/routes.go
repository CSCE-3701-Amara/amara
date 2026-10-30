package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/platform/config"
	"github.com/CSCE-3701-Amara/amara/internal/platform/database"
)

// RouteRegistrar mounts a module's HTTP handlers under /v1. Each module applies
// its own authentication and authorization middleware to protected routes.
type RouteRegistrar func(*gin.RouterGroup)

func newRouter(cfg config.HTTPConfig, logger *slog.Logger, dbSvc database.Service, registrars ...RouteRegistrar) (*gin.Engine, error) {
	router := gin.New()
	router.HandleMethodNotAllowed = true
	router.RedirectTrailingSlash = false
	if err := router.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		return nil, fmt.Errorf("configure trusted proxies: %w", err)
	}

	router.Use(requestID(), requestLogger(logger), recoverPanic(logger), corsMiddleware(cfg.AllowedOrigins))
	registerRoutes(router, dbSvc, registrars...)
	return router, nil
}

func registerRoutes(router *gin.Engine, dbSvc database.Service, registrars ...RouteRegistrar) {
	// /health/live reports process liveness without depending on the database.
	// /health/ready checks database connectivity when the service is wired.
	// Keep /health as a liveness alias for existing infrastructure probes.
	liveHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
	health := router.Group("/health")
	health.GET("", liveHandler)
	health.GET("/live", liveHandler)
	health.GET("/ready", func(c *gin.Context) {
		if dbSvc == nil {
			// The server can also run in configurations that do not use a database.
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}
		health := dbSvc.Health(c.Request.Context())
		status := http.StatusOK
		if health.Status != "up" {
			status = http.StatusServiceUnavailable
		}
		body := gin.H{
			"status": health.Status,
			"database": gin.H{
				"status": health.Status,
				"stats":  health.Stats,
			},
		}
		if health.Status != "up" {
			body["error"] = health.Error
		}
		c.JSON(status, body)
	})

	v1 := router.Group("/v1")
	for _, register := range registrars {
		register(v1)
	}
	router.NoRoute(func(c *gin.Context) {
		writeError(c, http.StatusNotFound, "not_found", "The requested resource was not found.")
	})
	router.NoMethod(func(c *gin.Context) {
		writeError(c, http.StatusMethodNotAllowed, "method_not_allowed", "The request method is not allowed.")
	})
}

func writeError(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"error":      gin.H{"code": code, "message": message},
		"request_id": RequestID(c.Request.Context()),
	})
}

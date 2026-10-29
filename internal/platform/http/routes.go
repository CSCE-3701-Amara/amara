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
	// /health/live — liveness: is the process running? No DB check.
	// /health/ready — readiness: is the database reachable? Returns 503 when not.
	//
	// The legacy /health route is kept as a liveness alias so existing
	// infrastructure probes continue to work unchanged.
	health := router.Group("/health")
	{
		liveHandler := func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
		health.GET("", liveHandler)      // /health  (legacy liveness alias)
		health.GET("/live", liveHandler) // /health/live

		health.GET("/ready", func(c *gin.Context) {
			if dbSvc == nil {
				// No database service wired; treat as live-only deployment.
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
				return
			}
			hs := dbSvc.Health(c.Request.Context())
			status := http.StatusOK
			if hs.Status != "up" {
				status = http.StatusServiceUnavailable
			}
			// Expose status and stats; suppress the raw error string from
			// public responses to avoid leaking internal hostnames.
			body := gin.H{
				"status": hs.Status,
				"database": gin.H{
					"status": hs.Status,
					"stats": gin.H{
						"total_conns":        hs.Stats.TotalConns,
						"acquired_conns":     hs.Stats.AcquiredConns,
						"idle_conns":         hs.Stats.IdleConns,
						"constructing_conns": hs.Stats.ConstructingConns,
						"max_conns":          hs.Stats.MaxConns,
						"acquire_count":      hs.Stats.AcquireCount,
						"empty_acquire_count": hs.Stats.EmptyAcquireCount,
					},
				},
			}
			if hs.Status != "up" {
				// Include a generic, non-sensitive reason.
				body["error"] = hs.Error
			}
			c.JSON(status, body)
		})
	}

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


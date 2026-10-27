package http

import (
	"context"
	"crypto/rand"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type requestIDKey struct{}

// RequestID returns the server-generated correlation ID from a request context.
func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

func requestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := rand.Text()
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), requestIDKey{}, id))
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

func requestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		defer func() {
			route := c.FullPath()
			if route == "" {
				route = "unmatched"
			}
			level := slog.LevelInfo
			if c.Writer.Status() >= http.StatusInternalServerError {
				level = slog.LevelError
			}
			// Route templates avoid logging identifiers, query tokens, or bodies.
			logger.Log(c.Request.Context(), level, "HTTP request",
				"request_id", RequestID(c.Request.Context()),
				"method", c.Request.Method,
				"route", route,
				"status", c.Writer.Status(),
				"duration_ms", float64(time.Since(start).Microseconds())/1000,
			)
		}()
		c.Next()
	}
}

func recoverPanic(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				if recovered == http.ErrAbortHandler {
					panic(recovered)
				}
				// Panic values may contain credentials or personal data.
				logger.ErrorContext(c.Request.Context(), "HTTP handler panic recovered",
					"request_id", RequestID(c.Request.Context()),
					"stack", string(debug.Stack()),
				)
				if c.Writer.Written() {
					// A partial response cannot safely become a JSON error.
					panic(http.ErrAbortHandler)
				}
				writeError(c, http.StatusInternalServerError, "internal_error", "An internal server error occurred.")
			}
		}()
		c.Next()
	}
}

func corsMiddleware(origins []string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			for _, allowed := range origins {
				if origin == allowed {
					return true
				}
			}
			return false
		},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders: []string{"X-Request-ID"},
		MaxAge:        12 * time.Hour,
	})
}

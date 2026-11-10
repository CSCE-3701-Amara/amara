package http

import (
	"net/http"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
	"github.com/gin-gonic/gin"
)

// RegisterCurrentSessionRoutes reads the existing session without issuing or extending it.
func RegisterCurrentSessionRoutes(group *gin.RouterGroup, validator ports.SessionValidator) {
	group.GET("/auth/sessions/current", func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Next()
	}, AuthenticatedAuth(validator), func(c *gin.Context) {
		identity, ok := Identity(c)
		if !ok {
			authenticationRequired(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"user":    gin.H{"id": identity.AccountID, "role": loginRole(identity.AccountType)},
			"session": gin.H{"expiresAt": identity.ExpiresAt.UTC().Format(time.RFC3339)},
		})
	})
}

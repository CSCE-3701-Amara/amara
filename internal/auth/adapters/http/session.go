package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

const identityContextKey = "amara.identity"

func CandidateAuth(validator ports.SessionValidator) gin.HandlerFunc {
	return sessionAuth(validator, domain.AccountTypeCandidate)
}

func sessionAuth(validator ports.SessionValidator, requiredType domain.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		if validator == nil {
			authenticationRequired(c)
			return
		}
		cookie, err := c.Request.Cookie(sessionCookieName)
		if err != nil || cookie.Value == "" {
			authenticationRequired(c)
			return
		}
		identity, err := validator.ValidateSession(c.Request.Context(), cookie.Value)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidSession) {
				authenticationRequired(c)
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "authentication_unavailable", "message": "Authentication is temporarily unavailable."},
			})
			return
		}
		if identity.AccountType != requiredType {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{"code": "forbidden", "message": "This account cannot access the requested resource."},
			})
			return
		}
		c.Set(identityContextKey, identity)
		c.Next()
	}
}

func authenticationRequired(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": gin.H{"code": "authentication_required", "message": "Authentication is required."},
	})
}

func Identity(c *gin.Context) (ports.Identity, bool) {
	value, exists := c.Get(identityContextKey)
	if !exists {
		return ports.Identity{}, false
	}
	identity, ok := value.(ports.Identity)
	return identity, ok
}

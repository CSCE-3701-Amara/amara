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

func AuthenticatedAuth(validator ports.SessionValidator) gin.HandlerFunc {
	return sessionAuth(validator, "")
}

func CompanyAuth(validator ports.SessionValidator, writeError func(*gin.Context, int, string, string)) gin.HandlerFunc {
	return sessionAuthWithErrors(validator, domain.AccountTypeCompany, writeError)
}

func sessionAuth(validator ports.SessionValidator, requiredType domain.AccountType) gin.HandlerFunc {
	return sessionAuthWithErrors(validator, requiredType, func(c *gin.Context, status int, code, message string) {
		c.AbortWithStatusJSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
	})
}

func sessionAuthWithErrors(validator ports.SessionValidator, requiredType domain.AccountType,
	writeError func(*gin.Context, int, string, string)) gin.HandlerFunc {
	return func(c *gin.Context) {
		if validator == nil {
			writeError(c, http.StatusUnauthorized, "authentication_required", "Authentication is required.")
			return
		}
		cookie, err := c.Request.Cookie(sessionCookieName)
		if err != nil || cookie.Value == "" {
			writeError(c, http.StatusUnauthorized, "authentication_required", "Authentication is required.")
			return
		}
		identity, err := validator.ValidateSession(c.Request.Context(), cookie.Value)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidSession) {
				writeError(c, http.StatusUnauthorized, "authentication_required", "Authentication is required.")
				return
			}
			writeError(c, http.StatusInternalServerError, "authentication_unavailable", "Authentication is temporarily unavailable.")
			return
		}
		if requiredType != "" && identity.AccountType != requiredType {
			writeError(c, http.StatusForbidden, "forbidden", "This account cannot access the requested resource.")
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

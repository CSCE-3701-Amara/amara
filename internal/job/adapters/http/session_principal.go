package http

import (
	"github.com/gin-gonic/gin"

	authhttp "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/http"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

// principalFromSession adapts AUTH's validated identity to the JOB handlers.
// Company operations are scoped by the company ID, not an interviewer ID.
func principalFromSession(c *gin.Context) {
	identity, ok := authhttp.Identity(c)
	if !ok {
		return // Handlers still fail closed; legacy handler tests supply their own principal.
	}
	id := identity.AccountID
	if identity.AccountType == domain.AccountTypeCompany {
		id = identity.CompanyID
	}
	c.Set(AuthenticatedAccountIDKey, id)
	c.Set(AuthenticatedAccountTypeKey, string(identity.AccountType))
}

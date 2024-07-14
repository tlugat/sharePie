package middleware

import (
	"net/http"
	"sharePie-api/internal/auth"
	"sharePie-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

func IsGranted(allowedRoles ...utils.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := auth.GetUserFromContext(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}
		for _, role := range allowedRoles {
			if role == user.Role {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
	}
}

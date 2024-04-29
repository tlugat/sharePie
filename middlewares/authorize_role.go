package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sharePie-api/helpers"
	"sharePie-api/utils"
)

func AuthorizeRole(allowedRoles ...utils.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := helpers.GetUserFromContext(c)
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

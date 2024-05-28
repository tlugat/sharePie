package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sharePie-api/internal/models"
)

func GetUserFromContext(c *gin.Context) (models.User, bool) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return models.User{}, false
	}

	user, ok := userInterface.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User information could not be retrieved"})
		return models.User{}, false
	}

	return user, true
}

package helpers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sharePie-api/models"
)

func GetUserFromContext(c *gin.Context) (*models.User, bool) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return nil, false
	}

	user, ok := userInterface.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User information could not be retrieved"})
		return nil, false
	}

	return &user, true
}

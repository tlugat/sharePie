package middleware

import (
	"net/http"
	"sharePie-api/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func IsEventAuthor(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			return
		}

		authUser, ok := user.(models.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User type assertion failed"})
			c.Abort()
			return
		}

		eventID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
			c.Abort()
			return
		}

		var event models.Event
		if err := db.First(&event, eventID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			c.Abort()
			return
		}

		if !IsUserEventAuthor(authUser, event) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is not the author of the event"})
			c.Abort()
			return
		}
		c.Next()
	}
}

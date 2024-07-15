package middleware

import (
	"net/http"
	"sharePie-api/internal/auth"
	"sharePie-api/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func EventIsUserPartOfEvent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := auth.GetUserFromContext(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		eventID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
			return
		}

		var event models.Event
		if err := db.First(&event, eventID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}

		result := db.Preload("Users").Where("id = ?", eventID).First(&event)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Fetching user event error"})
			return
		}

		var users []models.User
		for _, user := range event.Users {
			users = append(users, user)
		}

		event.Users = users

		if !IsUserPartOfEvent(user, event) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User is not part of the event"})
			return
		}
		c.Next()
	}
}

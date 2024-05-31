package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"sharePie-api/internal/models"

	"strconv"
)

func IsEventActive(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventID, err := strconv.Atoi(c.Param("eventId"))
		if err != nil {
			fmt.Println("Invalid event ID")
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

		if event.State != models.EventStateActive {
			c.JSON(http.StatusForbidden, gin.H{"error": "Event is archived"})
			c.Abort()
			return
		}

		c.Next()
	}
}

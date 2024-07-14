package middleware

import (
	"net/http"
	"sharePie-api/internal/auth"
	"sharePie-api/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func IsPartOfEvent_event(db *gorm.DB) gin.HandlerFunc {
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

		if !IsUserPartOfEvent(user, event) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User is not part of the event"})
			return
		}
		c.Next()
	}
}

func IsPartOfEvent_expense(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := auth.GetUserFromContext(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		expenseID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
			return
		}

		var expense models.Expense
		if err := db.First(&expense, expenseID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
			return
		}

		var event models.Event
		if err := db.First(&event, expense.EventID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}

		if !IsUserPartOfEvent(user, event) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User is not part of the event"})
			return
		}
		c.Next()
	}
}

func IsPartOfEvent_refund(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := auth.GetUserFromContext(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		refundID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid refund ID"})
			return
		}

		var refund models.Refund
		if err := db.First(&refund, refundID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Refund not found"})
			return
		}

		var event models.Event
		if err := db.First(&event, refund.EventID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}

		if !IsUserPartOfEvent(user, event) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User is not part of the event"})
			return
		}
		c.Next()
	}
}

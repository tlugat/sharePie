package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"sharePie-api/internal/auth"
	"sharePie-api/internal/models"
	"strconv"
)

func ExpenseIsUserPartOfEvent(db *gorm.DB) gin.HandlerFunc {
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

		result := db.Preload("Users").Where("id = ?", expense.EventID).First(&event)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Fetching user event error"})
			return
		}

		var users []models.User
		for _, user := range event.Users {
			fmt.Printf("user of event ==> %v\n", user.ID)
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

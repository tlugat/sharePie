package middleware

import (
	"sharePie-api/internal/models"
	"sharePie-api/pkg/utils"
)

func IsUserEventAuthor(user models.User, event models.Event) bool {
	if user.Role == utils.AdminRole {
		return true
	}
	if event.AuthorID != user.ID {
		return false
	}
	return true
}

func IsUserPartOfEvent(user models.User, event models.Event) bool {
	if IsUserEventAuthor(user, event) {
		return true
	}

	for _, participant := range event.Users {
		if participant.ID == user.ID {
			return true
		}
	}
	return false
}

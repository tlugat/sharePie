package types

import "sharePie-api/internal/models"

type IAchievementRepository interface {
	Find() ([]models.Achievement, error)
	FindOne(id uint) (models.Achievement, error)
	Create(achievement models.Achievement) (models.Achievement, error)
	Update(achievement models.Achievement) (models.Achievement, error)
	Delete(id uint) error
}

type IAchievementService interface {
	Find() ([]models.Achievement, error)
	FindOne(id uint) (models.Achievement, error)
	Create(input CreateAchievementInput) (models.Achievement, error)
	Update(id uint, input UpdateAchievementInput) (models.Achievement, error)
	Delete(id uint) error
}

type CreateAchievementInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Points      int    `json:"points" binding:"required"`
	Condition   string `json:"condition" binding:"required"`
}

type UpdateAchievementInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Points      int    `json:"points"`
	Condition   string `json:"condition"`
}

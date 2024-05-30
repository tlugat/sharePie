package types

import "sharePie-api/internal/models"

type ICategoryService interface {
	Find() ([]models.Category, error)
	FindOne(id uint) (models.Category, error)
	Create(input CreateCategoryInput) (models.Category, error)
	Update(id uint, input UpdateCategoryInput) (models.Category, error)
	Delete(id uint) error
}

type ICategoryRepository interface {
	Find() ([]models.Category, error)
	FindOne(id uint) (models.Category, error)
	Create(category models.Category) (models.Category, error)
	Update(category models.Category) (models.Category, error)
	Delete(id uint) error
}

type CreateCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCategoryInput struct {
	Name string `json:"name"`
}

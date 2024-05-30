package types

import "sharePie-api/internal/models"

type ITagService interface {
	Find() ([]models.Tag, error)
	FindOne(id uint) (models.Tag, error)
	Create(input CreateTagInput) (models.Tag, error)
	Update(id uint, input UpdateTagInput) (models.Tag, error)
	Delete(id uint) error
}

type ITagRepository interface {
	Find() ([]models.Tag, error)
	FindOne(id uint) (models.Tag, error)
	Create(tag models.Tag) (models.Tag, error)
	Update(tag models.Tag) (models.Tag, error)
	Delete(id uint) error
}

type CreateTagInput struct {
	Name string `json:"name" binding:"required"`
}

type UpdateTagInput struct {
	Name string `json:"title"`
}

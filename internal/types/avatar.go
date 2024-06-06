package types

import "sharePie-api/internal/models"

type IAvatarRepository interface {
	Find() ([]models.Avatar, error)
	FindOne(id uint) (models.Avatar, error)
	Create(avatar models.Avatar) (models.Avatar, error)
	Update(avatar models.Avatar) (models.Avatar, error)
	Delete(id uint) error
}

type IAvatarService interface {
	Find() ([]models.Avatar, error)
	FindOne(id uint) (models.Avatar, error)
	Create(input CreateAvatarInput) (models.Avatar, error)
	Update(id uint, input UpdateAvatarInput) (models.Avatar, error)
	Delete(id uint) error
}

type CreateAvatarInput struct {
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required"`
}

type UpdateAvatarInput struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

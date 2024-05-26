package types

import "sharePie-api/internal/models"

type IEventRepository interface {
	Find() ([]models.Event, error)
	FindOne(id uint) (models.Event, error)
	Create(event models.Event) (models.Event, error)
	Update(event models.Event) (models.Event, error)
	Delete(id uint) error
	FindOneByCode(code string) (models.Event, error)
	FindUsers(id uint) ([]models.User, error)
}

type IEventService interface {
	Find() ([]models.Event, error)
	FindOne(id uint) (models.Event, error)
	Create(input CreateEventInput, user models.User) (models.Event, error)
	Update(id uint, input UpdateEventInput) (models.Event, error)
	Delete(id uint) error
	GetUsers(id uint) ([]models.User, error)
	AddUser(code string, user models.User) error
}

type CreateEventInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Category    uint    `json:"category"`
	Image       string  `json:"image"`
	Goal        float64 `json:"goal"`
}

type UpdateEventInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    uint    `json:"category"`
	Image       string  `json:"image"`
	Goal        float64 `json:"goal"`
}

type JoinEventInput struct {
	Code string `json:"code" binding:"required"`
}

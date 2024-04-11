package services

import (
	"go-project/models"
	"go-project/repositories"
)

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

type IEventService interface {
	Find() ([]models.Event, error)
	FindOne(id uint) (models.Event, error)
	Create(input CreateEventInput, user models.User) (models.Event, error)
	Update(id uint, input UpdateEventInput) (models.Event, error)
	Delete(id uint) error
}

type EventService struct {
	Repository         repositories.IEventRepository
	CategoryRepository repositories.ICategoryRepository
	UserRepository     repositories.IUserRepository
}

func NewEventService(
	repository repositories.IEventRepository,
	categoryRepository repositories.ICategoryRepository,
	userRepository repositories.IUserRepository) IEventService {
	return &EventService{
		Repository:         repository,
		CategoryRepository: categoryRepository,
		UserRepository:     userRepository,
	}
}

func (service *EventService) Find() ([]models.Event, error) {
	events, err := service.Repository.Find()
	if err != nil {
		return nil, err
	}

	// Fetch category and author for each event
	for i, event := range events {
		category, err := service.CategoryRepository.FindOne(event.CategoryID)
		author, err := service.UserRepository.FindOneById(event.AuthorID)

		if err != nil {
			return nil, err
		}
		events[i].Category = category
		events[i].Author = author
	}

	return events, nil
}

func (service *EventService) FindOne(id uint) (models.Event, error) {
	event, err := service.Repository.FindOne(id)

	if err != nil {
		return models.Event{}, err
	}
	// Fetch event's category
	category, err := service.CategoryRepository.FindOne(event.CategoryID)
	// Fetch event's author
	author, err := service.UserRepository.FindOneById(event.AuthorID)

	if err != nil {
		return models.Event{}, err
	}

	event.Category = category
	event.Author = author

	return event, nil

}

func (service *EventService) Create(input CreateEventInput, user models.User) (models.Event, error) {
	event := models.Event{
		Name:        input.Name,
		Description: input.Description,
		CategoryID:  input.Category,
		Image:       input.Image,
		Goal:        input.Goal,
		AuthorID:    user.ID,
	}

	return service.Repository.Create(event)
}

func (service *EventService) Update(id uint, input UpdateEventInput) (models.Event, error) {
	event, err := service.Repository.FindOne(id)

	if err != nil {
		return models.Event{}, err
	}

	if input.Name != "" {
		event.Name = input.Name
	}
	if input.Description != "" {
		event.Description = input.Description
	}
	if input.Category != 0 {
		event.CategoryID = input.Category
	}
	if input.Image != "" {
		event.Image = input.Image
	}
	if input.Goal != 0 {
		event.Goal = input.Goal
	}

	return service.Repository.Update(event)
}

func (service *EventService) Delete(id uint) error {
	return service.Repository.Delete(id)
}

package services

import (
	"sharePie-api/models"
	"sharePie-api/repositories"
	"sharePie-api/utils"
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
	GetUsers(id uint) ([]models.User, error)
	AddUser(code string, user models.User) error
}

type EventService struct {
	Repository         repositories.IEventRepository
	CategoryRepository repositories.ICategoryRepository
	UserRepository     repositories.IUserRepository
	ExpenseRepository  repositories.IExpenseRepository
}

func NewEventService(
	repository repositories.IEventRepository,
	categoryRepository repositories.ICategoryRepository,
	userRepository repositories.IUserRepository,
	expenseRepository repositories.IExpenseRepository,
) IEventService {
	return &EventService{
		Repository:         repository,
		CategoryRepository: categoryRepository,
		UserRepository:     userRepository,
		ExpenseRepository:  expenseRepository,
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
		Code:        utils.GenerateInvitationCode(6),
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

func (service *EventService) GetUsers(id uint) ([]models.User, error) {
	users, err := service.UserRepository.FindByEventId(id)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (service *EventService) AddUser(code string, user models.User) error {
	event, err := service.Repository.FindOneByCode(code)
	if err != nil {
		return err
	}

	users, err := service.UserRepository.FindByEventId(event.ID)
	if err != nil {
		return err
	}
	isUserAlreadyInEvent := false

	for _, u := range users {
		if u.ID == user.ID {
			isUserAlreadyInEvent = true
			break
		}
	}

	if isUserAlreadyInEvent {
		return nil
	}

	event.Users = append(users, user)

	_, err = service.Repository.Update(event)

	return err
}

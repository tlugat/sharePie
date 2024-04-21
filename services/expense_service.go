package services

import (
	"sharePie-api/models"
	"sharePie-api/repositories"
)

type CreateExpenseInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Tag         uint    `json:"tag"`
	Payer       uint    `json:"payer" binding:"required"`
	Image       string  `json:"image"`
	Amount      float64 `json:"amount" binding:"required"`
	Event       uint    `json:"event" binding:"required"`
	Users       []uint  `json:"users" binding:"required"`
}

type UpdateExpenseInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Tag         uint    `json:"category"`
	Payer       uint    `json:"payer"`
	Image       string  `json:"image"`
	Amount      float64 `json:"amount"`
	Users       []uint  `json:"users"`
}

type IExpenseService interface {
	Find() ([]models.Expense, error)
	FindOne(id uint) (models.Expense, error)
	Create(input CreateExpenseInput, user models.User) (models.Expense, error)
	Update(id uint, input UpdateExpenseInput) (models.Expense, error)
	Delete(id uint) error
}

type ExpenseService struct {
	Repository     repositories.IExpenseRepository
	TagRepository  repositories.ITagRepository
	UserRepository repositories.IUserRepository
}

func NewExpenseService(
	repository repositories.IExpenseRepository,
	categoryRepository repositories.ITagRepository,
	userRepository repositories.IUserRepository) IExpenseService {
	return &ExpenseService{
		Repository:     repository,
		TagRepository:  categoryRepository,
		UserRepository: userRepository,
	}
}

func (service *ExpenseService) Find() ([]models.Expense, error) {
	expenses, err := service.Repository.Find()
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func (service *ExpenseService) FindOne(id uint) (models.Expense, error) {
	expense, err := service.Repository.FindOne(id)

	if err != nil {
		return models.Expense{}, err
	}

	return expense, nil
}

func (service *ExpenseService) Create(input CreateExpenseInput, user models.User) (models.Expense, error) {
	expense := models.Expense{
		Name:        input.Name,
		Description: input.Description,
		TagID:       input.Tag,
		PayerID:     input.Payer,
		Image:       input.Image,
		Amount:      input.Amount,
		AuthorID:    user.ID,
		EventID:     input.Event,
	}

	var users []models.User
	if err := service.UserRepository.FindByIds(input.Users, &users); err != nil {
		return models.Expense{}, err
	}

	expense.Users = users

	return service.Repository.Create(expense)
}

func (service *ExpenseService) Update(id uint, input UpdateExpenseInput) (models.Expense, error) {
	expense, err := service.Repository.FindOne(id)

	if err != nil {
		return models.Expense{}, err
	}

	if input.Name != "" {
		expense.Name = input.Name
	}
	if input.Description != "" {
		expense.Description = input.Description
	}
	if input.Tag != 0 {
		expense.TagID = input.Tag
	}
	if input.Payer != 0 {
		expense.PayerID = input.Payer
	}
	if input.Image != "" {
		expense.Image = input.Image
	}
	if input.Amount != 0 {
		expense.Amount = input.Amount
	}
	if input.Users != nil {
		var users []models.User
		if err := service.UserRepository.FindByIds(input.Users, &users); err != nil {
			return models.Expense{}, err
		}
		expense.Users = users
	}

	return service.Repository.Update(expense)
}

func (service *ExpenseService) Delete(id uint) error {
	return service.Repository.Delete(id)
}

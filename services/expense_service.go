package services

import (
	"sharePie-api/models"
	"sharePie-api/repositories"
)

type CreateExpenseInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Tag         uint    `json:"tag"`
	Image       string  `json:"image"`
	Amount      float64 `json:"amount" binding:"required"`
	Event       uint    `json:"event" binding:"required"`
}

type UpdateExpenseInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Tag         uint    `json:"category"`
	Image       string  `json:"image"`
	Amount      float64 `json:"amount" binding:"required"`
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

	// Fetch category and author for each expense
	for i, expense := range expenses {
		category, err := service.TagRepository.FindOne(expense.TagID)
		author, err := service.UserRepository.FindOneById(expense.AuthorID)

		if err != nil {
			return nil, err
		}
		expenses[i].Tag = category
		expenses[i].Author = author
	}

	return expenses, nil
}

func (service *ExpenseService) FindOne(id uint) (models.Expense, error) {
	expense, err := service.Repository.FindOne(id)

	if err != nil {
		return models.Expense{}, err
	}
	// Fetch expense's category
	tag, err := service.TagRepository.FindOne(expense.TagID)
	// Fetch expense's author
	author, err := service.UserRepository.FindOneById(expense.AuthorID)

	if err != nil {
		return models.Expense{}, err
	}

	expense.Tag = tag
	expense.Author = author

	return expense, nil

}

func (service *ExpenseService) Create(input CreateExpenseInput, user models.User) (models.Expense, error) {
	expense := models.Expense{
		Name:        input.Name,
		Description: input.Description,
		TagID:       input.Tag,
		Image:       input.Image,
		Amount:      input.Amount,
		AuthorID:    user.ID,
		EventID:     input.Event,
	}

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
	if input.Image != "" {
		expense.Image = input.Image
	}
	if input.Amount != 0 {
		expense.Amount = input.Amount
	}

	return service.Repository.Update(expense)
}

func (service *ExpenseService) Delete(id uint) error {
	return service.Repository.Delete(id)
}

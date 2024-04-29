package expense

import (
	models2 "sharePie-api/internal/models"
	"sharePie-api/internal/tag"
	"sharePie-api/internal/user"
	"sharePie-api/pkg/config/thirdparty/cloudinary"
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
	Find() ([]models2.Expense, error)
	FindOne(id uint) (models2.Expense, error)
	Create(input CreateExpenseInput, user models2.User) (models2.Expense, error)
	Update(id uint, input UpdateExpenseInput) (models2.Expense, error)
	Delete(id uint) error
}

type Service struct {
	Repository     IExpenseRepository
	TagRepository  tag.ITagRepository
	UserRepository user.IUserRepository
}

func NewService(
	repository IExpenseRepository,
	categoryRepository tag.ITagRepository,
	userRepository user.IUserRepository) IExpenseService {
	return &Service{
		Repository:     repository,
		TagRepository:  categoryRepository,
		UserRepository: userRepository,
	}
}

func (service *Service) Find() ([]models2.Expense, error) {
	expenses, err := service.Repository.Find()
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func (service *Service) FindOne(id uint) (models2.Expense, error) {
	expense, err := service.Repository.FindOne(id)

	if err != nil {
		return models2.Expense{}, err
	}

	return expense, nil
}

func (service *Service) Create(input CreateExpenseInput, user models2.User) (models2.Expense, error) {
	image, err := cloudinary.UploadImage(input.Image, "Events")
	if err != nil {
		return models2.Expense{}, err
	}
	expense := models2.Expense{
		Name:        input.Name,
		Description: input.Description,
		TagID:       input.Tag,
		PayerID:     input.Payer,
		Image:       image,
		Amount:      input.Amount,
		AuthorID:    user.ID,
		EventID:     input.Event,
	}

	var users []models2.User
	if err := service.UserRepository.FindByIds(input.Users, &users); err != nil {
		return models2.Expense{}, err
	}

	expense.Users = users

	return service.Repository.Create(expense)
}

func (service *Service) Update(id uint, input UpdateExpenseInput) (models2.Expense, error) {
	expense, err := service.Repository.FindOne(id)

	if err != nil {
		return models2.Expense{}, err
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
		image, err := cloudinary.UploadImage(input.Image, "Events")
		if err != nil {
			return models2.Expense{}, err
		}
		expense.Image = image
	}
	if input.Amount != 0 {
		expense.Amount = input.Amount
	}
	if input.Users != nil {
		var users []models2.User
		if err := service.UserRepository.FindByIds(input.Users, &users); err != nil {
			return models2.Expense{}, err
		}
		expense.Users = users
	}

	return service.Repository.Update(expense)
}

func (service *Service) Delete(id uint) error {
	return service.Repository.Delete(id)
}

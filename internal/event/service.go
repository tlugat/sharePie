package event

import (
	"math/rand"
	"sharePie-api/internal/category"
	models2 "sharePie-api/internal/models"
	"sharePie-api/internal/types"
	"sharePie-api/internal/user"
	"sharePie-api/pkg/config/thirdparty/cloudinary"
	"strings"
	"time"
)

type Service struct {
	Repository         types.IEventRepository
	CategoryRepository category.ICategoryRepository
	UserRepository     user.IUserRepository
	ExpenseRepository  types.IExpenseRepository
}

func NewService(
	repository types.IEventRepository,
	categoryRepository category.ICategoryRepository,
	userRepository user.IUserRepository,
	expenseRepository types.IExpenseRepository,
) types.IEventService {
	return &Service{
		Repository:         repository,
		CategoryRepository: categoryRepository,
		UserRepository:     userRepository,
		ExpenseRepository:  expenseRepository,
	}
}

func (service *Service) Find() ([]models2.Event, error) {
	events, err := service.Repository.Find()
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (service *Service) FindOne(id uint) (models2.Event, error) {
	event, err := service.Repository.FindOne(id)

	if err != nil {
		return models2.Event{}, err
	}

	return event, nil

}

func (service *Service) Create(input types.CreateEventInput, user models2.User) (models2.Event, error) {
	event := models2.Event{
		Name:        input.Name,
		Description: input.Description,
		CategoryID:  input.Category,
		Goal:        input.Goal,
		AuthorID:    user.ID,
		Code:        generateInvitationCode(6),
		Users:       []models2.User{user},
	}
	if input.Image != "" {
		image, err := cloudinary.UploadImage(input.Image, "Events")
		if err != nil {
			return models2.Event{}, err
		}
		event.Image = image
	}

	return service.Repository.Create(event)
}

func (service *Service) Update(id uint, input types.UpdateEventInput) (models2.Event, error) {
	event, err := service.Repository.FindOne(id)

	if err != nil {
		return models2.Event{}, err
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
		image, err := cloudinary.UploadImage(input.Image, "Events")
		if err != nil {
			return models2.Event{}, err
		}
		event.Image = image
	}
	if input.Goal != 0 {
		event.Goal = input.Goal
	}

	return service.Repository.Update(event)
}

func (service *Service) Delete(id uint) error {
	return service.Repository.Delete(id)
}

func (service *Service) GetUsers(id uint) ([]models2.User, error) {
	users, err := service.UserRepository.FindByEventId(id)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (service *Service) AddUser(code string, user models2.User) error {
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

func generateInvitationCode(length int) string {
	var chars = "ABCDEFGHJKLMNPQRSTUVWXYZ123456789"
	var result strings.Builder
	result.Grow(length)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		index := rand.Intn(len(chars))
		result.WriteByte(chars[index])
	}

	return result.String()
}

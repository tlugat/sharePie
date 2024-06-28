package event

import (
	"errors"
	"math/rand"
	models2 "sharePie-api/internal/models"
	"sharePie-api/internal/types"
	"sharePie-api/pkg/config/thirdparty/cloudinary"
	"strings"
	"time"
)

type Service struct {
	Repository         types.IEventRepository
	CategoryRepository types.ICategoryRepository
	UserRepository     types.IUserRepository
	ExpenseRepository  types.IExpenseRepository
}

func NewService(
	repository types.IEventRepository,
	categoryRepository types.ICategoryRepository,
	userRepository types.IUserRepository,
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
	category, err := service.CategoryRepository.FindOne(input.Category)
	if err != nil {
		return models2.Event{}, err
	}

	event := models2.Event{
		Name:        input.Name,
		Description: input.Description,
		CategoryID:  input.Category,
		Category:    category,
		Goal:        input.Goal,
		AuthorID:    user.ID,
		Author:      user,
		Code:        generateInvitationCode(6),
		Users:       []models2.User{user},
		State:       models2.EventStateActive,
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

	if input.Users != nil {
		users, err := service.UserRepository.FindByIds(input.Users)
		if err != nil {
			return models2.Event{}, err
		}
		err = service.Repository.RemoveUsers(event)
		if err != nil {
			return models2.Event{}, err
		}
		event.Users = users

	}

	return service.Repository.Update(event)
}

func (service *Service) UpdateState(id uint, input types.UpdateEventStateInput) (models2.Event, error) {
	event, err := service.Repository.FindOne(id)
	if err != nil {
		return models2.Event{}, err
	}

	if err := input.State.IsValid(); err != nil {
		return models2.Event{}, err
	}

	event.State = input.State

	return service.Repository.Update(event)
}

func (service *Service) Delete(id uint) error {
	return service.Repository.Delete(id)
}

func (service *Service) GetUsers(eventID uint) ([]types.UserWithExpenses, error) {
	users, err := service.UserRepository.FindByEventId(eventID)
	if err != nil {
		return nil, err
	}

	var usersWithExpenses []types.UserWithExpenses

	for _, user := range users {
		expenses, err := service.ExpenseRepository.FindByPayerUserIdAndEventId(user.ID, eventID)
		if err != nil {
			return nil, err
		}

		totalExpenses := 0.0
		for _, expense := range expenses {
			for _, payer := range expense.Payers {
				if payer.UserID == user.ID {
					totalExpenses += payer.Amount
				}
			}
		}

		userWithExpenses := types.UserWithExpenses{
			User:          user,
			ExpenseCount:  len(expenses),
			TotalExpenses: totalExpenses,
		}
		usersWithExpenses = append(usersWithExpenses, userWithExpenses)
	}

	return usersWithExpenses, nil
}

func (service *Service) AddUser(code string, user models2.User) (models2.Event, error) {
	event, err := service.Repository.FindOneByCode(code)
	if err != nil {
		return models2.Event{}, err
	}

	if event.State != models2.EventStateActive {
		return models2.Event{}, errors.New("Event is not active")
	}

	users, err := service.UserRepository.FindByEventId(event.ID)
	if err != nil {
		return models2.Event{}, err
	}
	isUserAlreadyInEvent := false

	for _, u := range users {
		if u.ID == user.ID {
			isUserAlreadyInEvent = true
			break
		}
	}

	if isUserAlreadyInEvent {
		return models2.Event{}, types.NewConflictError("User is already in the event")
	}

	event.Users = append(users, user)

	_, err = service.Repository.Update(event)
	if err != nil {
		return models2.Event{}, err
	}

	return event, nil
}

func (service *Service) GetBalances(event models2.Event) ([]models2.Balance, error) {
	balances, err := service.Repository.FindBalances(event)
	if err != nil {
		return nil, err
	}
	return balances, nil
}

func (service *Service) GetTransactions(event models2.Event) ([]models2.Transaction, error) {
	transactions, err := service.Repository.FindTransactions(event)
	if err != nil {
		return nil, err
	}
	return transactions, nil
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

func (service *Service) FindExpenses(id uint) ([]models2.Expense, error) {
	expenses, err := service.ExpenseRepository.FindByEventId(id)
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

func (service *Service) FindByUser(id uint) ([]models2.Event, error) {
	events, err := service.Repository.FindByUser(id)
	if err != nil {
		return nil, err
	}
	return events, nil
}

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

func (service *Service) CreateBalances(event models2.Event) ([]models2.Balance, error) {
	expenses, err := service.ExpenseRepository.FindByEventId(event.ID)
	if err != nil {
		return nil, err
	}

	eventUsers, err := service.Repository.FindUsers(event.ID)

	if err != nil {
		return nil, err
	}

	userBalances := make(map[uint]float64)
	userDetails := make(map[uint]models2.User)

	for _, eventUser := range eventUsers {
		userBalances[eventUser.ID] = 0
		userDetails[eventUser.ID] = eventUser
	}

	var total float64

	for _, eventExpense := range expenses {
		total += eventExpense.Amount
		for _, participant := range eventExpense.Participants {
			userBalances[participant.UserID] -= participant.Amount
		}
		for _, payer := range eventExpense.Payers {
			userBalances[payer.UserID] += payer.Amount
		}
	}

	var balances []models2.Balance
	for userID, userBalance := range userBalances {
		balance := models2.Balance{
			UserID:  userID,
			User:    userDetails[userID],
			Amount:  userBalance,
			EventID: event.ID,
		}
		balances = append(balances, balance)
	}
	err = service.Repository.DeleteBalances(event)
	if err != nil {
		return nil, err
	}
	err = service.Repository.CreateBalances(balances)
	if err != nil {
		return nil, err
	}

	return balances, nil
}

func (service *Service) CreateTransactions(event models2.Event) ([]models2.Transaction, error) {
	balances, err := service.GetBalances(event)
	if err != nil {
		return nil, err
	}

	var creditors []models2.Balance
	var debtors []models2.Balance

	for _, balance := range balances {
		if balance.Amount > 0 {
			creditors = append(creditors, balance)
		} else if balance.Amount < 0 {
			debtors = append(debtors, balance)
		}
	}

	var transactions []models2.Transaction

	for len(creditors) > 0 && len(debtors) > 0 {
		creditor := creditors[0]
		debtor := debtors[0]

		amount := min(creditor.Amount, -debtor.Amount)

		transaction := models2.Transaction{
			FromUserID: debtor.UserID,
			From:       debtor.User,
			ToUserID:   creditor.UserID,
			To:         creditor.User,
			EventID:    event.ID,
			Completed:  false,
			Amount:     amount,
		}
		transactions = append(transactions, transaction)

		creditor.Amount -= amount
		debtor.Amount += amount

		if creditor.Amount == 0 {
			creditors = creditors[1:]
		} else {
			creditors[0] = creditor
		}

		if debtor.Amount == 0 {
			debtors = debtors[1:]
		} else {
			debtors[0] = debtor
		}
	}

	err = service.Repository.DeleteTransactions(event)
	if err != nil {
		return nil, err
	}
	err = service.Repository.CreateTransactions(transactions)
	if err != nil {
		return nil, err
	}

	return transactions, nil
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
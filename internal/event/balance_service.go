package event

import (
	"sharePie-api/internal/expense"
	models2 "sharePie-api/internal/models"
	"sharePie-api/internal/user"
)

type BalanceSummaryUser struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
}

type BalanceSummary struct {
	User   BalanceSummaryUser `json:"user"`
	Amount float64            `json:"amount"`
}

type IEventBalanceService interface {
	GetBalanceSummary(event models2.Event) ([]BalanceSummary, error)
}

type BalanceService struct {
	EventRepository   IEventRepository
	ExpenseRepository expense.IExpenseRepository
	UserRepository    user.IUserRepository
}

func NewBalanceService(
	eventRepository IEventRepository,
	expenseRepository expense.IExpenseRepository,
	userRepository user.IUserRepository,
) IEventBalanceService {
	return &BalanceService{
		EventRepository:   eventRepository,
		ExpenseRepository: expenseRepository,
		UserRepository:    userRepository,
	}
}

func (service *BalanceService) GetBalanceSummary(event models2.Event) ([]BalanceSummary, error) {
	expenses, err := service.ExpenseRepository.FindByEventId(event.ID)
	if err != nil {
		return nil, err
	}

	eventUsers, err := service.UserRepository.FindByEventId(event.ID)

	if err != nil {
		return nil, err
	}

	userBalances := make(map[uint]float64)
	userDetails := make(map[uint]BalanceSummaryUser)

	for _, eventUser := range eventUsers {
		userBalances[eventUser.ID] = 0
		userDetails[eventUser.ID] = BalanceSummaryUser{
			ID:        eventUser.ID,
			FirstName: eventUser.FirstName,
			LastName:  eventUser.LastName,
			Username:  eventUser.Username,
		}
	}

	var total float64

	for _, eventExpense := range expenses {
		total += eventExpense.Amount
		isUserConcerned := isUserConcerned(eventExpense.Users, eventExpense.PayerID)

		if !isUserConcerned {
			userBalances[eventExpense.PayerID] += eventExpense.Amount
		}
		for _, eventUser := range eventExpense.Users {
			if eventExpense.PayerID == eventUser.ID {
				userBalances[eventUser.ID] += eventExpense.Amount - (eventExpense.Amount / float64(len(eventExpense.Users)))
				continue
			}
			userBalances[eventUser.ID] -= eventExpense.Amount / float64(len(eventExpense.Users))
		}
	}

	var summaries []BalanceSummary
	for userID, balance := range userBalances {
		summary := BalanceSummary{
			User:   userDetails[userID],
			Amount: balance,
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func isUserConcerned(users []models2.User, userID uint) bool {
	for _, eventUser := range users {
		if eventUser.ID == userID {
			return true
		}
	}
	return false
}

package services

import (
	"sharePie-api/models"
	"sharePie-api/repositories"
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
	GetBalanceSummary(event models.Event) ([]BalanceSummary, error)
}

type EventBalanceService struct {
	EventRepository   repositories.IEventRepository
	ExpenseRepository repositories.IExpenseRepository
	UserRepository    repositories.IUserRepository
}

func NewEventBalanceService(
	eventRepository repositories.IEventRepository,
	expenseRepository repositories.IExpenseRepository,
	userRepository repositories.IUserRepository,
) IEventBalanceService {
	return &EventBalanceService{
		EventRepository:   eventRepository,
		ExpenseRepository: expenseRepository,
		UserRepository:    userRepository,
	}
}

func (service *EventBalanceService) GetBalanceSummary(event models.Event) ([]BalanceSummary, error) {
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

	for _, user := range eventUsers {
		userBalances[user.ID] = 0
		userDetails[user.ID] = BalanceSummaryUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
		}
	}

	var total float64

	for _, expense := range expenses {
		total += expense.Amount
		isUserConcerned := isUserConcerned(expense.Users, expense.PayerID)

		if !isUserConcerned {
			userBalances[expense.PayerID] += expense.Amount
		}
		for _, user := range expense.Users {
			if expense.PayerID == user.ID {
				userBalances[user.ID] += expense.Amount - (expense.Amount / float64(len(expense.Users)))
				continue
			}
			userBalances[user.ID] -= expense.Amount / float64(len(expense.Users))
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

func isUserConcerned(users []models.User, userID uint) bool {
	for _, user := range users {
		if user.ID == userID {
			return true
		}
	}
	return false
}

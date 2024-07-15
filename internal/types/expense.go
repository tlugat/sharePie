package types

import "sharePie-api/internal/models"

type IExpenseRepository interface {
	Find() ([]models.Expense, error)
	FindByEventId(id uint) ([]models.Expense, error)
	FindByUserIdAndEventId(userID uint, eventID uint) ([]models.Expense, error)
	FindByPayerUserIdAndEventId(userID uint, eventID uint) ([]models.Expense, error)
	FindOne(id uint) (models.Expense, error)
	Create(expense models.Expense) (models.Expense, error)
	Update(expense models.Expense) (models.Expense, error)
	Delete(id uint) error
}

type IExpenseService interface {
	Find() ([]models.Expense, error)
	FindOne(id uint) (models.Expense, error)
	Create(input CreateExpenseInput, user models.User) (models.Expense, error)
	Update(id uint, input UpdateExpenseInput) (models.Expense, error)
	Delete(id uint) error
	FindByEventId(id uint) ([]models.Expense, error)
}

type ParticipantInput struct {
	Id     uint    `json:"id"`
	Amount float64 `json:"amount"`
}
type PayerInput struct {
	Id     uint    `json:"id"`
	Amount float64 `json:"amount"`
}

type CreateExpenseInput struct {
	Name         string             `json:"name" binding:"required"`
	Description  string             `json:"description" binding:"required"`
	Tag          uint               `json:"tag"`
	Image        string             `json:"image"`
	Amount       float64            `json:"amount" binding:"required"`
	Event        uint               `json:"event" binding:"required"`
	Participants []ParticipantInput `json:"participants" binding:"required"`
	Payers       []PayerInput       `json:"payers" binding:"required"`
}

type UpdateExpenseInput struct {
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Tag          uint               `json:"tag"`
	Image        string             `json:"image"`
	Participants []ParticipantInput `json:"participants"`
	Payers       []PayerInput       `json:"payers"`
	Amount       float64            `json:"amount"`
	ID           uint               `json:"id"`
}

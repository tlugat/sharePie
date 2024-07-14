package types

import (
	"errors"
	"sharePie-api/internal/models"
)

type IEventRepository interface {
	Find() ([]models.Event, error)
	FindOne(id uint) (models.Event, error)
	Create(event models.Event) (models.Event, error)
	Update(event models.Event) (models.Event, error)
	UpdateUsers(event models.Event) ([]models.User, error)
	Delete(id uint) error
	FindOneByCode(code string) (models.Event, error)
	FindUsers(id uint) ([]models.User, error)
	CreateBalances(balances []models.Balance) error
	CreateTransactions(transactions []models.Transaction) error
	FindBalances(event models.Event) ([]models.Balance, error)
	FindTransactions(event models.Event) ([]models.Transaction, error)
	DeleteBalances(event models.Event) error
	DeleteTransactions(event models.Event) error
	RemoveUsers(event models.Event) error
	FindByUser(id uint) ([]models.Event, error)
}

type IEventService interface {
	Find() ([]models.Event, error)
	FindOne(id uint) (models.Event, error)
	Create(input CreateEventInput, user models.User) (models.Event, error)
	Update(id uint, input UpdateEventInput) (models.Event, error)
	UpdateState(id uint, input UpdateEventStateInput) (models.Event, error)
	Delete(id uint) error
	GetUsers(id uint) ([]UserWithExpenses, error)
	AddUser(code string, user models.User) (models.Event, error)
	GetBalances(event models.Event) ([]models.Balance, error)
	GetTransactions(event models.Event) ([]models.Transaction, error)
	FindExpenses(id uint) ([]models.Expense, error)
	FindByUser(id uint) ([]models.Event, error)
	CreateBalances(event models.Event) ([]models.Balance, error)
	CreateTransactions(event models.Event, balances []models.Balance) ([]models.Transaction, error)
}

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
	Users       []uint  `json:"users"`
}

type JoinEventInput struct {
	Code string `json:"code" binding:"required"`
}

type UpdateEventStateInput struct {
	State models.EventState `json:"state" binding:"required"`
}

type Transaction struct {
	From   models.User `json:"from"`
	To     models.User `json:"to"`
	Amount float64     `json:"amount"`
}

var ErrConflict = errors.New("conflict")

type ConflictError struct {
	Message string
}

type UserWithExpenses struct {
	models.User
	ExpenseCount  int     `json:"expense_count"`
	TotalExpenses float64 `json:"total_expenses"`
}

func (e *ConflictError) Error() string {
	return e.Message
}

func NewConflictError(message string) error {
	return &ConflictError{Message: message}
}

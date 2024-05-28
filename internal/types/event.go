package types

import "sharePie-api/internal/models"

type IEventRepository interface {
	Find() ([]models.Event, error)
	FindOne(id uint) (models.Event, error)
	Create(event models.Event) (models.Event, error)
	Update(event models.Event) (models.Event, error)
	Delete(id uint) error
	FindOneByCode(code string) (models.Event, error)
	FindUsers(id uint) ([]models.User, error)
	CreateBalances(balances []models.Balance) error
	CreateTransactions(transactions []models.Transaction) error
	FindBalances(event models.Event) ([]models.Balance, error)
	FindTransactions(event models.Event) ([]models.Transaction, error)
	DeleteBalances(event models.Event) error
	DeleteTransactions(event models.Event) error
}

type IEventService interface {
	Find() ([]models.Event, error)
	FindOne(id uint) (models.Event, error)
	Create(input CreateEventInput, user models.User) (models.Event, error)
	Update(id uint, input UpdateEventInput) (models.Event, error)
	Delete(id uint) error
	GetUsers(id uint) ([]models.User, error)
	AddUser(code string, user models.User) error
	CreateBalances(event models.Event) ([]models.Balance, error)
	CreateTransactions(event models.Event) ([]models.Transaction, error)
	GetBalances(event models.Event) ([]models.Balance, error)
	GetTransactions(event models.Event) ([]models.Transaction, error)
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

type BalanceUser struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
}

type Transaction struct {
	From   BalanceUser `json:"from"`
	To     BalanceUser `json:"to"`
	Amount float64     `json:"amount"`
}

package models

import (
	"gorm.io/gorm"
)

type Expense struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	TagId       uint
	Tag         Tag     `json:"tag"`
	Users       []User  `json:"users" gorm:"many2many:expense_users;"`
	Image       string  `json:"image"`
	Amount      float64 `string:"cost"`
	EventID     uint    `json:"eventId"`
}

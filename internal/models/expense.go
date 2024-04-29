package models

import (
	"gorm.io/gorm"
)

type Expense struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	TagID       uint    `json:"-" `
	Tag         Tag     `json:"tag" gorm:"foreignKey:TagID"`
	AuthorID    uint    `json:"-" `
	Author      User    `json:"author" gorm:"foreignKey:AuthorID"`
	PayerID     uint    `json:"-" `
	Payer       User    `json:"payer" gorm:"foreignKey:PayerID"`
	Users       []User  `json:"users" gorm:"many2many:expense_users;"`
	Image       string  `json:"image"`
	Amount      float64 `string:"cost"`
	EventID     uint    `json:"eventId"`
}

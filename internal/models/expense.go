package models

import (
	"gorm.io/gorm"
	"time"
)

type Expense struct {
	gorm.Model
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	TagID        uint          `json:"-" `
	Tag          Tag           `json:"tag" gorm:"foreignKey:TagID"`
	Image        string        `json:"image"`
	AuthorID     uint          `json:"-" `
	Author       User          `json:"author" gorm:"foreignKey:AuthorID"`
	Payers       []Payer       `json:"payers" gorm:"foreignKey:ExpenseID"`
	Participants []Participant `json:"participants" gorm:"foreignKey:ExpenseID"`
	Amount       float64       `json:"amount"`
	EventID      uint          `json:"eventId"`
	Date         time.Time     `json:"date" time_format:"2006-01-02T15:04:05Z07:00"`
}

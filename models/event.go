package models

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	Name        string   `json:"name"`
	Description string   `json:"description"`
	AuthorID    uint     `json:"-"`
	Author      User     `json:"author" gorm:"-"`
	CategoryID  uint     `json:"-"`
	Category    Category `json:"category" gorm:"-"`
	Users       []User   `json:"users" gorm:"many2many:event_users;"`
	Image       string   `json:"image"`
	Goal        float64  `json:"goal"`
	Expenses    []Expense
}

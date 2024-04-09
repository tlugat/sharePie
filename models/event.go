package models

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      uint   `json:"author_id"`
	CategoryID  uint
	Category    uint   `json:"category"`
	Users       []User `json:"users" gorm:"many2many:event_users;"`
	Image       string `json:"image"`
	Goal        int    `json:"goal"`
	Expenses    []Expense
}

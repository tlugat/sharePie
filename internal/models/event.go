package models

import (
	"errors"
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name        string     `json:"name"`
	Description string     `json:"description"`
	AuthorID    uint       `json:"-"`
	Author      User       `json:"author" gorm:"foreignKey:AuthorID"`
	CategoryID  uint       `json:"-"`
	Category    Category   `json:"category" gorm:"foreignKey:CategoryID"`
	Users       []User     `json:"-" gorm:"many2many:event_users;"`
	Image       string     `json:"image"`
	Goal        float64    `json:"goal"`
	Expenses    []Expense  `json:"-"`
	Code        string     `json:"code" gorm:"unique"`
	State       EventState `json:"state"`
}

type EventState string

const (
	EventStateActive   EventState = "active"
	EventStateArchived EventState = "archived"
)

func (es EventState) IsValid() error {
	switch es {
	case EventStateActive, EventStateArchived:
		return nil
	}
	return errors.New("invalid event state")
}

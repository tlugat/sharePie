package models

import (
	"gorm.io/gorm"
	"sharePie-api/utils"
)

type User struct {
	gorm.Model
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Email     string     `gorm:"unique"`
	Password  string     `json:"-"`
	Role      utils.Role `json:"role"`
	Events    []Event    `json:"-" gorm:"foreignKey:AuthorID"`
}

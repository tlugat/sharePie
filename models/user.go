package models

import (
	"go-project/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `gorm:"unique"`
	Password  string
	Role      utils.Role `json:"role"`
}

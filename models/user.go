package models

import (
	"go-project/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
	Role     utils.Role `json:"role"`
}

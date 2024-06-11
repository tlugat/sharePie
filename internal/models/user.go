package models

import (
	"gorm.io/gorm"
	"sharePie-api/pkg/constants"
)

type User struct {
	gorm.Model
	Username     string         `json:"username" gorm:"unique"`
	Email        string         `gorm:"unique"`
	Password     string         `json:"-"`
	Role         constants.Role `json:"role"`
	Events       []Event        `json:"-" gorm:"foreignKey:AuthorID"`
	Achievements []Achievement  `json:"-" gorm:"many2many:user_achievements;"`
	AvatarID     *uint          `json:"-"`
	Avatar       *Avatar        `json:"avatar" gorm:"foreignKey:AvatarID"`
}

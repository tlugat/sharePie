package models

import (
	"sharePie-api/pkg/utils"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username      string        `json:"username" gorm:"unique"`
	Email         string        `json:"email" gorm:"unique"`
	Password      string        `json:"-"`
	Role          utils.Role    `json:"role"`
	Events        []Event       `json:"-" gorm:"foreignKey:AuthorID"`
	Achievements  []Achievement `json:"-" gorm:"many2many:user_achievements;"`
	AvatarID      uint          `json:"-"`
	Avatar        Avatar        `json:"avatar" gorm:"foreignKey:AvatarID"`
	FirebaseToken *string       `json:"firebaseToken" gorm:"unique"`
}

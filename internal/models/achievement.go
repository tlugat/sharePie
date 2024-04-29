package models

import "gorm.io/gorm"

type Achievement struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique"`
	Description string `json:"description"`
	Points      int    `json:"points"`
	Condition   string `json:"condition"`
}

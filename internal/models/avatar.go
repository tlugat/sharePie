package models

import "gorm.io/gorm"

type Avatar struct {
	gorm.Model
	Name string `json:"name" gorm:"unique"`
	URL  string `json:"url"`
}

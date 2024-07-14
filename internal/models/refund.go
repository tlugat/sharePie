package models

import (
	"gorm.io/gorm"
	"time"
)

type Refund struct {
	gorm.Model
	FromUserID uint      `json:"-" gorm:"column:from;not null"`
	From       User      `json:"from" gorm:"foreignKey:FromUserID;references:ID;"`
	ToUserID   uint      `json:"-" gorm:"column:to;not null"`
	To         User      `json:"to" gorm:"foreignKey:ToUserID;references:ID;"`
	Amount     float64   `json:"amount"`
	EventID    uint      `json:"eventId" gorm:"not null"`
	Event      Event     `json:"-" gorm:"foreignKey:EventID;references:ID"`
	AuthorID   uint      `json:"-" gorm:"column:author;not null"`
	Author     User      `json:"author" gorm:"foreignKey:AuthorID;references:ID"`
	Date       time.Time `json:"date" time_format:"2006-01-02T15:04:05Z07:00"`
}

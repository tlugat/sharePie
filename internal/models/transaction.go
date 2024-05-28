package models

type Transaction struct {
	ID         uint    `json:"id" gorm:"primaryKey"`
	FromUserID uint    `json:"-" gorm:"column:from;not null"`
	From       User    `json:"from" gorm:"foreignKey:FromUserID;references:ID;"`
	ToUserID   uint    `json:"-" gorm:"column:to;not null"`
	To         User    `json:"to" gorm:"foreignKey:ToUserID;references:ID;"`
	Amount     float64 `json:"amount"`
	EventID    uint    `json:"event_id" gorm:"not null"`
	Event      Event   `json:"-" gorm:"foreignKey:EventID;references:ID"`
	Completed  bool    `json:"completed"`
}

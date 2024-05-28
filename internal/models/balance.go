package models

type Balance struct {
	ID      uint    `json:"id" gorm:"primaryKey"`
	UserID  uint    `json:"-"  gorm:"foreignKey:UserID"`
	User    User    `json:"user"`
	Amount  float64 `json:"amount"`
	EventID uint    `json:"event_id" gorm:"foreignKey:EventID"`
}

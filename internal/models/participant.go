package models

type Participant struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	UserID    uint    `json:"-"`
	User      User    `json:"user" gorm:"foreignKey:UserID"`
	Amount    float64 `json:"amount"`
	ExpenseID uint    `json:"-"`
}

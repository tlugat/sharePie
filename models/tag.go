package models

type Tag struct {
	ID       uint      `json:"id" gorm:"primary_key"`
	Name     string    `json:"name"`
	Expenses []Expense `json:"-" gorm:"foreignKey:TagID"`
}

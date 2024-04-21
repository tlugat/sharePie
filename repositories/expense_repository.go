package repositories

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sharePie-api/models"
)

type IExpenseRepository interface {
	Find() ([]models.Expense, error)
	FindByEventId(id uint) ([]models.Expense, error)
	FindOne(id uint) (models.Expense, error)
	Create(expense models.Expense) (models.Expense, error)
	Update(expense models.Expense) (models.Expense, error)
	Delete(id uint) error
}

type ExpenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) IExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Find() ([]models.Expense, error) {
	var expenses []models.Expense
	result := r.db.Preload(clause.Associations).Find(&expenses)
	return expenses, result.Error
}

func (r *ExpenseRepository) FindByEventId(id uint) ([]models.Expense, error) {
	var expenses []models.Expense
	result := r.db.Preload(clause.Associations).Where("event_id = ?", id).Find(&expenses)
	return expenses, result.Error
}

func (r *ExpenseRepository) FindOne(id uint) (models.Expense, error) {
	var expense models.Expense
	result := r.db.Preload(clause.Associations).Where("id = ?", id).First(&expense)
	return expense, result.Error
}

func (r *ExpenseRepository) Create(expense models.Expense) (models.Expense, error) {
	result := r.db.Create(&expense)
	return expense, result.Error
}

func (r *ExpenseRepository) Update(expense models.Expense) (models.Expense, error) {
	result := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&expense)
	return expense, result.Error
}

func (r *ExpenseRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Expense{}, id)
	return result.Error
}

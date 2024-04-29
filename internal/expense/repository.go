package expense

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sharePie-api/internal/models"
)

type IExpenseRepository interface {
	Find() ([]models.Expense, error)
	FindByEventId(id uint) ([]models.Expense, error)
	FindOne(id uint) (models.Expense, error)
	Create(expense models.Expense) (models.Expense, error)
	Update(expense models.Expense) (models.Expense, error)
	Delete(id uint) error
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) IExpenseRepository {
	return &Repository{db: db}
}

func (r *Repository) Find() ([]models.Expense, error) {
	var expenses []models.Expense
	result := r.db.Preload(clause.Associations).Find(&expenses)
	return expenses, result.Error
}

func (r *Repository) FindByEventId(id uint) ([]models.Expense, error) {
	var expenses []models.Expense
	result := r.db.Preload(clause.Associations).Where("event_id = ?", id).Find(&expenses)
	return expenses, result.Error
}

func (r *Repository) FindOne(id uint) (models.Expense, error) {
	var expense models.Expense
	result := r.db.Preload(clause.Associations).Where("id = ?", id).First(&expense)
	return expense, result.Error
}

func (r *Repository) Create(expense models.Expense) (models.Expense, error) {
	result := r.db.Create(&expense)
	return expense, result.Error
}

func (r *Repository) Update(expense models.Expense) (models.Expense, error) {
	result := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&expense)
	return expense, result.Error
}

func (r *Repository) Delete(id uint) error {
	result := r.db.Delete(&models.Expense{}, id)
	return result.Error
}

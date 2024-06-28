package expense

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sharePie-api/internal/models"
	"sharePie-api/internal/types"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) types.IExpenseRepository {
	return &Repository{db: db}
}

func (r *Repository) Find() ([]models.Expense, error) {
	var expenses []models.Expense
	result := r.db.Preload(clause.Associations).Preload("Author.Avatar").Preload("Participants.User.Avatar").
		Preload("Payers.User.Avatar").Find(&expenses)
	return expenses, result.Error
}

func (r *Repository) FindByEventId(id uint) ([]models.Expense, error) {
	var expenses []models.Expense
	result := r.db.Preload(clause.Associations).Preload("Author.Avatar").Preload("Participants.User.Avatar").
		Preload("Payers.User.Avatar").Where("event_id = ?", id).Find(&expenses)
	return expenses, result.Error
}

func (r *Repository) FindByUserIdAndEventId(userID uint, eventID uint) ([]models.Expense, error) {
	var expenses []models.Expense
	result := r.db.Where("author_id = ? AND event_id = ?", userID, eventID).Find(&expenses)
	return expenses, result.Error
}

func (r *Repository) FindByPayerUserIdAndEventId(userID uint, eventID uint) ([]models.Expense, error) {
	var expenses []models.Expense
	err := r.db.Joins("JOIN payers ON payers.expense_id = expenses.id").
		Where("payers.user_id = ? AND expenses.event_id = ?", userID, eventID).
		Preload("Payers").
		Preload("Participants").
		Find(&expenses).Error
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *Repository) FindOne(id uint) (models.Expense, error) {
	var expense models.Expense
	result := r.db.Preload("Author.Avatar").Preload(clause.Associations).Preload("Participants.User.Avatar").
		Preload("Payers.User.Avatar").Where("id = ?", id).First(&expense)
	return expense, result.Error
}

func (r *Repository) Create(expense models.Expense) (models.Expense, error) {
	result := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&expense)
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

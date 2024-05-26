package payer

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sharePie-api/internal/models"
)

type IPayerRepository interface {
	Find() ([]models.Payer, error)
	FindByExpenseId(id uint) ([]models.Payer, error)
	FindOne(id uint) (models.Payer, error)
	Create(payer models.Payer) (models.Payer, error)
	Update(payer models.Payer) (models.Payer, error)
	Delete(id uint) error
	DeleteByExpenseID(id uint) error
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) IPayerRepository {
	return &Repository{db: db}
}

func (r *Repository) Find() ([]models.Payer, error) {
	var payers []models.Payer
	result := r.db.Preload(clause.Associations).Find(&payers)
	return payers, result.Error
}

func (r *Repository) FindByExpenseId(id uint) ([]models.Payer, error) {
	var payers []models.Payer
	result := r.db.Preload(clause.Associations).Where("expense_id = ?", id).Find(&payers)
	return payers, result.Error
}

func (r *Repository) FindOne(id uint) (models.Payer, error) {
	var payer models.Payer
	result := r.db.Preload(clause.Associations).Where("id = ?", id).First(&payer)
	return payer, result.Error
}

func (r *Repository) Create(payer models.Payer) (models.Payer, error) {
	result := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&payer)
	return payer, result.Error
}

func (r *Repository) Update(payer models.Payer) (models.Payer, error) {
	result := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&payer)
	return payer, result.Error
}

func (r *Repository) Delete(id uint) error {
	result := r.db.Delete(&models.Payer{}, id)
	return result.Error
}

func (r *Repository) DeleteByExpenseID(id uint) error {
	err := r.db.Where("expense_id = ?", id).Delete(&[]models.Payer{}).Error
	return err
}

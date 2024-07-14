package refund

import (
	"gorm.io/gorm"
	"sharePie-api/internal/models"
	"sharePie-api/internal/types"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) types.IRefundRepository {
	return &Repository{db: db}
}

func (r *Repository) Find() ([]models.Refund, error) {
	var refunds []models.Refund
	result := r.db.Preload("From.Avatar").Preload("To.Avatar").Preload("Author.Avatar").Find(&refunds)
	return refunds, result.Error
}

func (r *Repository) FindOne(id uint) (models.Refund, error) {
	var refund models.Refund
	result := r.db.Preload("From.Avatar").Preload("To.Avatar").Preload("Author.Avatar").Where("id = ?", id).First(&refund)
	return refund, result.Error
}

func (r *Repository) FindByEventId(eventId uint) ([]models.Refund, error) {
	var refunds []models.Refund
	result := r.db.Preload("From.Avatar").Preload("To.Avatar").Preload("Author.Avatar").Where("event_id = ?", eventId).Find(&refunds)
	return refunds, result.Error
}

func (r *Repository) Create(refund models.Refund) (models.Refund, error) {
	result := r.db.Create(&refund)
	return refund, result.Error
}

func (r *Repository) Update(refund models.Refund) (models.Refund, error) {
	result := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&refund)
	return refund, result.Error
}

func (r *Repository) Delete(id uint) error {
	result := r.db.Delete(&models.Refund{}, id)
	return result.Error
}

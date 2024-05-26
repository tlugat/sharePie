package participant

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sharePie-api/internal/models"
)

type IParticipantRepository interface {
	Find() ([]models.Participant, error)
	FindByExpenseId(id uint) ([]models.Participant, error)
	FindOne(id uint) (models.Participant, error)
	Create(participant models.Participant) (models.Participant, error)
	Update(participant models.Participant) (models.Participant, error)
	Delete(id uint) error
	DeleteByExpenseId(id uint) error
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) IParticipantRepository {
	return &Repository{db: db}
}

func (r *Repository) Find() ([]models.Participant, error) {
	var participants []models.Participant
	result := r.db.Preload(clause.Associations).Find(&participants)
	return participants, result.Error
}

func (r *Repository) FindByExpenseId(id uint) ([]models.Participant, error) {
	var participants []models.Participant
	result := r.db.Preload(clause.Associations).Where("expense_id = ?", id).Find(&participants)
	return participants, result.Error
}

func (r *Repository) FindOne(id uint) (models.Participant, error) {
	var participant models.Participant
	result := r.db.Preload(clause.Associations).Where("id = ?", id).First(&participant)
	return participant, result.Error
}

func (r *Repository) Create(participant models.Participant) (models.Participant, error) {
	result := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&participant)
	return participant, result.Error
}

func (r *Repository) Update(participant models.Participant) (models.Participant, error) {
	result := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&participant)
	return participant, result.Error
}

func (r *Repository) Delete(id uint) error {
	result := r.db.Delete(&models.Participant{}, id)
	return result.Error
}

func (r *Repository) DeleteByExpenseId(id uint) error {
	result := r.db.Delete(&models.Participant{}, "expense_id = ?", id)
	return result.Error
}

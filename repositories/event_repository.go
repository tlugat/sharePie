package repositories

import (
	"go-project/models"
	"gorm.io/gorm"
)

type IEventRepository interface {
	Find() ([]models.Event, error)
	FindOne(id uint) (models.Event, error)
	Create(event models.Event) (models.Event, error)
	Update(event models.Event) (models.Event, error)
	Delete(id uint) error
}

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) IEventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Find() ([]models.Event, error) {
	var events []models.Event
	result := r.db.Find(&events)
	return events, result.Error
}

func (r *EventRepository) FindOne(id uint) (models.Event, error) {
	var event models.Event
	result := r.db.Where("id = ?", id).First(&event)
	return event, result.Error
}

func (r *EventRepository) Create(event models.Event) (models.Event, error) {
	result := r.db.Create(&event)
	return event, result.Error
}

func (r *EventRepository) Update(event models.Event) (models.Event, error) {
	result := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&event)
	return event, result.Error
}

func (r *EventRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Event{}, id)
	return result.Error
}

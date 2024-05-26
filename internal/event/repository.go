package event

import (
	"gorm.io/gorm"
	"sharePie-api/internal/models"
	"sharePie-api/internal/types"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) types.IEventRepository {
	return &Repository{db: db}
}

func (r *Repository) Find() ([]models.Event, error) {
	var events []models.Event
	result := r.db.Preload("Category").Preload("Author").Find(&events)
	return events, result.Error
}

func (r *Repository) FindOne(id uint) (models.Event, error) {
	var event models.Event
	result := r.db.Preload("Category").Preload("Author").Where("id = ?", id).First(&event)
	return event, result.Error
}

func (r *Repository) Create(event models.Event) (models.Event, error) {
	result := r.db.Create(&event)
	return event, result.Error
}

func (r *Repository) Update(event models.Event) (models.Event, error) {
	result := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&event)
	return event, result.Error
}

func (r *Repository) Delete(id uint) error {
	result := r.db.Delete(&models.Event{}, id)
	return result.Error
}

func (r *Repository) FindOneByCode(code string) (models.Event, error) {
	var event models.Event
	result := r.db.Where("code = ?", code).First(&event)
	return event, result.Error
}

func (r *Repository) FindUsers(id uint) ([]models.User, error) {
	var event models.Event
	result := r.db.Preload("Users").Where("id = ?", id).First(&event)
	if result.Error != nil {
		return nil, result.Error
	}

	var users []models.User
	for _, user := range event.Users {
		users = append(users, user)
	}
	return users, nil
}

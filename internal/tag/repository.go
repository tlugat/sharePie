package tag

import (
	"gorm.io/gorm"
	"sharePie-api/internal/models"
)

type ITagRepository interface {
	Find() ([]models.Tag, error)
	FindOne(id uint) (models.Tag, error)
	Create(tag models.Tag) (models.Tag, error)
	Update(tag models.Tag) (models.Tag, error)
	Delete(id uint) error
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) ITagRepository {
	return &Repository{db: db}
}

func (r *Repository) Find() ([]models.Tag, error) {
	var tags []models.Tag
	result := r.db.Find(&tags)
	return tags, result.Error
}

func (r *Repository) FindOne(id uint) (models.Tag, error) {
	var tag models.Tag
	result := r.db.Where("id = ?", id).First(&tag)
	return tag, result.Error
}

func (r *Repository) Create(tag models.Tag) (models.Tag, error) {
	result := r.db.Create(&tag)
	return tag, result.Error
}

func (r *Repository) Update(tag models.Tag) (models.Tag, error) {
	result := r.db.Save(&tag)
	return tag, result.Error
}

func (r *Repository) Delete(id uint) error {
	result := r.db.Delete(&models.Tag{}, id)
	return result.Error
}
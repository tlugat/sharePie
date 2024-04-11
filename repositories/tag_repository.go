package repositories

import (
	"go-project/models"
	"gorm.io/gorm"
)

type ITagRepository interface {
	Find() ([]models.Tag, error)
	FindOne(id uint) (models.Tag, error)
	Create(tag models.Tag) (models.Tag, error)
	Update(tag models.Tag) (models.Tag, error)
	Delete(id uint) error
}

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) ITagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Find() ([]models.Tag, error) {
	var tags []models.Tag
	result := r.db.Find(&tags)
	return tags, result.Error
}

func (r *TagRepository) FindOne(id uint) (models.Tag, error) {
	var tag models.Tag
	result := r.db.Where("id = ?", id).First(&tag)
	return tag, result.Error
}

func (r *TagRepository) Create(tag models.Tag) (models.Tag, error) {
	result := r.db.Create(&tag)
	return tag, result.Error
}

func (r *TagRepository) Update(tag models.Tag) (models.Tag, error) {
	result := r.db.Save(&tag)
	return tag, result.Error
}

func (r *TagRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Tag{}, id)
	return result.Error
}

package category

import (
	"gorm.io/gorm"
	"sharePie-api/internal/models"
)

type ICategoryRepository interface {
	Find() ([]models.Category, error)
	FindOne(id uint) (models.Category, error)
	Create(category models.Category) (models.Category, error)
	Update(category models.Category) (models.Category, error)
	Delete(id uint) error
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) ICategoryRepository {
	return &Repository{db: db}
}

func (r *Repository) Find() ([]models.Category, error) {
	var categories []models.Category
	result := r.db.Find(&categories)
	return categories, result.Error
}

func (r *Repository) FindOne(id uint) (models.Category, error) {
	var category models.Category
	result := r.db.Where("id = ?", id).First(&category)
	return category, result.Error
}

func (r *Repository) Create(category models.Category) (models.Category, error) {
	result := r.db.Create(&category)
	return category, result.Error
}

func (r *Repository) Update(category models.Category) (models.Category, error) {
	result := r.db.Save(&category)
	return category, result.Error
}

func (r *Repository) Delete(id uint) error {
	result := r.db.Delete(&models.Category{}, id)
	return result.Error
}

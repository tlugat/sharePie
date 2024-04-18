package repositories

import (
	"gorm.io/gorm"
	"sharePie-api/models"
)

type ICategoryRepository interface {
	Find() ([]models.Category, error)
	FindOne(id uint) (models.Category, error)
	Create(category models.Category) (models.Category, error)
	Update(category models.Category) (models.Category, error)
	Delete(id uint) error
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) ICategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Find() ([]models.Category, error) {
	var categories []models.Category
	result := r.db.Find(&categories)
	return categories, result.Error
}

func (r *CategoryRepository) FindOne(id uint) (models.Category, error) {
	var category models.Category
	result := r.db.Where("id = ?", id).First(&category)
	return category, result.Error
}

func (r *CategoryRepository) Create(category models.Category) (models.Category, error) {
	result := r.db.Create(&category)
	return category, result.Error
}

func (r *CategoryRepository) Update(category models.Category) (models.Category, error) {
	result := r.db.Save(&category)
	return category, result.Error
}

func (r *CategoryRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Category{}, id)
	return result.Error
}

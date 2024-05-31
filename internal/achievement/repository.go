package achievement

import (
	"gorm.io/gorm"
	"sharePie-api/internal/models"
	"sharePie-api/internal/types"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) types.IAchievementRepository {
	return &Repository{db: db}
}

func (r *Repository) Find() ([]models.Achievement, error) {
	var achievements []models.Achievement
	result := r.db.Find(&achievements)
	return achievements, result.Error
}

func (r *Repository) FindOne(id uint) (models.Achievement, error) {
	var achievement models.Achievement
	result := r.db.Where("id = ?", id).First(&achievement)
	return achievement, result.Error
}

func (r *Repository) Create(achievement models.Achievement) (models.Achievement, error) {
	result := r.db.Create(&achievement)
	return achievement, result.Error
}

func (r *Repository) Update(achievement models.Achievement) (models.Achievement, error) {
	result := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&achievement)
	return achievement, result.Error
}

func (r *Repository) Delete(id uint) error {
	result := r.db.Delete(&models.Achievement{}, id)
	return result.Error
}

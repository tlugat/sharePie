package repositories

import (
	"gorm.io/gorm"
	"sharePie-api/models"
)

type IAchievementRepository interface {
	Find() ([]models.Achievement, error)
	FindOne(id uint) (models.Achievement, error)
	Create(achievement models.Achievement) (models.Achievement, error)
	Update(achievement models.Achievement) (models.Achievement, error)
	Delete(id uint) error
}

type AchievementRepository struct {
	db *gorm.DB
}

func NewAchievementRepository(db *gorm.DB) IAchievementRepository {
	return &AchievementRepository{db: db}
}

func (r *AchievementRepository) Find() ([]models.Achievement, error) {
	var achievements []models.Achievement
	result := r.db.Find(&achievements)
	return achievements, result.Error
}

func (r *AchievementRepository) FindOne(id uint) (models.Achievement, error) {
	var achievement models.Achievement
	result := r.db.Where("id = ?", id).First(&achievement)
	return achievement, result.Error
}

func (r *AchievementRepository) Create(achievement models.Achievement) (models.Achievement, error) {
	result := r.db.Create(&achievement)
	return achievement, result.Error
}

func (r *AchievementRepository) Update(achievement models.Achievement) (models.Achievement, error) {
	result := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&achievement)
	return achievement, result.Error
}

func (r *AchievementRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Achievement{}, id)
	return result.Error
}

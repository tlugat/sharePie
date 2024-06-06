package avatar

import (
	"gorm.io/gorm"
	"sharePie-api/internal/models"
	"sharePie-api/internal/types"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) types.IAvatarRepository {
	return &Repository{db: db}
}

func (r *Repository) Find() ([]models.Avatar, error) {
	var avatars []models.Avatar
	result := r.db.Find(&avatars)
	return avatars, result.Error
}

func (r *Repository) FindOne(id uint) (models.Avatar, error) {
	var avatar models.Avatar
	result := r.db.Where("id = ?", id).First(&avatar)
	return avatar, result.Error
}

func (r *Repository) Create(avatar models.Avatar) (models.Avatar, error) {
	result := r.db.Create(&avatar)
	return avatar, result.Error
}

func (r *Repository) Update(avatar models.Avatar) (models.Avatar, error) {
	result := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&avatar)
	return avatar, result.Error
}

func (r *Repository) Delete(id uint) error {
	result := r.db.Delete(&models.Avatar{}, id)
	return result.Error
}

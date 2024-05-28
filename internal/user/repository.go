package user

import (
	"gorm.io/gorm"
	"sharePie-api/internal/models"
)

type IUserRepository interface {
	Find() ([]models.User, error)
	FindByIds(ids []uint, users *[]models.User) error
	FindByEventId(eventId uint) ([]models.User, error)
	FindOneById(id uint) (models.User, error)
	FindOneByEmail(email string) (models.User, error)
	Create(user models.User) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(id uint) error
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) IUserRepository {
	return &Repository{db: db}
}

func (r *Repository) Find() ([]models.User, error) {
	var users []models.User
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *Repository) FindByIds(ids []uint, users *[]models.User) error {
	return r.db.Where("id IN ?", ids).Find(users).Error
}

func (r *Repository) FindByEventId(eventId uint) ([]models.User, error) {
	var users []models.User
	result := r.db.Joins("JOIN event_users ON users.id = event_users.user_id").Where("event_users.event_id = ?", eventId).Find(&users)
	return users, result.Error
}

func (r *Repository) FindOneById(id uint) (models.User, error) {
	var user models.User
	result := r.db.Where("id = ?", id).First(&user)
	return user, result.Error
}

func (r *Repository) FindOneByEmail(email string) (models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	return user, result.Error
}

func (r *Repository) Create(user models.User) (models.User, error) {
	result := r.db.Create(&user)
	return user, result.Error
}

func (r *Repository) Update(user models.User) (models.User, error) {
	result := r.db.Save(&user)
	return user, result.Error
}

func (r *Repository) Delete(id uint) error {
	result := r.db.Delete(&models.User{}, id)
	return result.Error
}
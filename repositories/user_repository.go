package repositories

import (
	"go-project/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Find() ([]models.User, error)
	FindOne(id int) (models.User, error)
	Create(user models.User) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(id int) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Find() ([]models.User, error) {
	var users []models.User
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *UserRepository) FindOne(id int) (models.User, error) {
	var user models.User
	result := r.db.Where("id = ?", id).First(&user)
	return user, result.Error
}

func (r *UserRepository) Create(user models.User) (models.User, error) {
	result := r.db.Create(&user)
	return user, result.Error
}

func (r *UserRepository) Update(user models.User) (models.User, error) {
	result := r.db.Save(&user)
	return user, result.Error
}

func (r *UserRepository) Delete(id int) error {
	result := r.db.Delete(&models.User{}, id)
	return result.Error
}

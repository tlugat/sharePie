package types

import "sharePie-api/internal/models"

type IUserService interface {
	Find() ([]models.User, error)
	FindOneById(id uint) (models.User, error)
	FindOneByEmail(email string) (models.User, error)
	Create(input CreateUserInput) (models.User, error)
	Update(id uint, input UpdateUserInput) (models.User, error)
	UpdateFirebaseToken(id uint, input UpdateUserFirebaseTokenInput) (models.User, error)
	Delete(id uint) error
}

type IUserRepository interface {
	Find() ([]models.User, error)
	FindByIds(ids []uint) ([]models.User, error)
	FindByEventId(eventId uint) ([]models.User, error)
	FindOneById(id uint) (models.User, error)
	FindOneByEmail(email string) (models.User, error)
	Create(user models.User) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(id uint) error
}

type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserInput struct {
	Username string `json:"username"`
	Avatar   uint   `json:"avatar"`
	Email    string `json:"email"`
}

type UpdateUserFirebaseTokenInput struct {
	FirebaseToken string `json:"firebaseToken" binding:"required"`
}

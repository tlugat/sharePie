package user

import (
	"sharePie-api/internal/models"
	"sharePie-api/internal/types"
	"sharePie-api/pkg/constants"
)

type Service struct {
	Repository types.IUserRepository
}

func NewService(repository types.IUserRepository) types.IUserService {
	return &Service{Repository: repository}
}

func (service *Service) Find() ([]models.User, error) {
	return service.Repository.Find()
}

func (service *Service) FindOneById(id uint) (models.User, error) {
	return service.Repository.FindOneById(id)
}

func (service *Service) FindOneByEmail(email string) (models.User, error) {
	return service.Repository.FindOneByEmail(email)
}

func (service *Service) Create(input types.CreateUserInput) (models.User, error) {
	user := models.User{Username: input.Username, Email: input.Email, Password: input.Password, Role: constants.UserRole}
	return service.Repository.Create(user)
}

func (service *Service) Update(id uint, input types.UpdateUserInput) (models.User, error) {
	user, err := service.Repository.FindOneById(id)

	if err != nil {
		return models.User{}, err
	}

	if input.Username != "" {
		user.Username = input.Username
	}

	return service.Repository.Update(user)
}

func (service *Service) Delete(id uint) error {
	return service.Repository.Delete(id)
}

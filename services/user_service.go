package services

import (
	"go-project/models"
	"go-project/repositories"
)

type CreateUserInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type IUserService interface {
	Find() ([]models.User, error)
	FindOne(id int) (models.User, error)
	Create(input CreateUserInput) (models.User, error)
	Update(id int, input UpdateUserInput) (models.User, error)
	Delete(id int) error
}

type UserService struct {
	Repository repositories.IUserRepository
}

func NewUserService(repository repositories.IUserRepository) IUserService {
	return &UserService{Repository: repository}
}

func (service *UserService) Find() ([]models.User, error) {
	return service.Repository.Find()
}

func (service *UserService) FindOne(id int) (models.User, error) {
	return service.Repository.FindOne(id)
}

func (service *UserService) Create(input CreateUserInput) (models.User, error) {
	user := models.User{FirstName: input.FirstName, LastName: input.LastName, Email: input.Email, Password: input.Password}
	return service.Repository.Create(user)
}

func (service *UserService) Update(id int, input UpdateUserInput) (models.User, error) {
	user, err := service.Repository.FindOne(id)

	if err != nil {
		return models.User{}, err
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}

	return service.Repository.Update(user)
}

func (service *UserService) Delete(id int) error {
	return service.Repository.Delete(id)
}

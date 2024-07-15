package user

import (
	"errors"
	"fmt"
	"sharePie-api/internal/models"
	"sharePie-api/internal/types"
	"sharePie-api/pkg/utils"
)

type Service struct {
	Repository       types.IUserRepository
	AvatarRepository types.IAvatarRepository
}

func NewService(repository types.IUserRepository, avatarRepository types.IAvatarRepository) types.IUserService {
	return &Service{Repository: repository, AvatarRepository: avatarRepository}
}

func (service *Service) Find() ([]models.User, error) {
	users, err := service.Repository.Find()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to find users: %v", err))
	}
	return users, nil
}

func (service *Service) FindOneById(id uint) (models.User, error) {
	user, err := service.Repository.FindOneById(id)
	if err != nil {
		return models.User{}, errors.New(fmt.Sprintf("failed to find user with id %d: %v", id, err))
	}
	return user, nil
}

func (service *Service) FindOneByEmail(email string) (models.User, error) {
	user, err := service.Repository.FindOneByEmail(email)
	if err != nil {
		return models.User{}, errors.New(fmt.Sprintf("failed to find user with email %s: %v", email, err))
	}
	return user, nil
}

func (service *Service) Create(input types.CreateUserInput) (models.User, error) {
	defaultAvatar, err := service.AvatarRepository.FindOne(25)
	if err != nil {
		return models.User{}, errors.New(fmt.Sprintf("failed to find default avatar: %v", err))
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		Role:     utils.UserRole,
		AvatarID: defaultAvatar.ID,
		Avatar:   defaultAvatar,
	}

	newUser, err := service.Repository.Create(user)
	if err != nil {
		return models.User{}, errors.New(fmt.Sprintf("failed to create user: %v", err))
	}
	return newUser, nil
}

func (service *Service) Update(id uint, input types.UpdateUserInput) (models.User, error) {
	user, err := service.Repository.FindOneById(id)
	if err != nil {
		return models.User{}, errors.New(fmt.Sprintf("failed to find user with id %d: %v", id, err))
	}

	if input.Username != "" {
		user.Username = input.Username
	}

	if input.Email != "" {
		user.Email = input.Email
	}

	if input.Avatar != 0 {
		avatar, err := service.AvatarRepository.FindOne(input.Avatar)
		if err != nil {
			return models.User{}, errors.New(fmt.Sprintf("failed to find avatar with id %d: %v", input.Avatar, err))
		}
		user.AvatarID = input.Avatar
		user.Avatar = avatar
	}

	updatedUser, err := service.Repository.Update(user)
	if err != nil {
		return models.User{}, errors.New(fmt.Sprintf("failed to update user with id %d: %v", id, err))
	}
	return updatedUser, nil
}

func (service *Service) UpdateFirebaseToken(id uint, input types.UpdateUserFirebaseTokenInput) (models.User, error) {
	user, err := service.Repository.FindOneById(id)
	if err != nil {
		return models.User{}, errors.New(fmt.Sprintf("failed to find user with id %d: %v", id, err))
	}

	if input.FirebaseToken != "" {
		user.FirebaseToken = &input.FirebaseToken
	}

	updatedUser, err := service.Repository.Update(user)
	if err != nil {
		return models.User{}, errors.New(fmt.Sprintf("failed to update Firebase token for user with id %d: %v", id, err))
	}
	return updatedUser, nil
}

func (service *Service) Delete(id uint) error {
	err := service.Repository.Delete(id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete user with id %d: %v", id, err))
	}
	return nil
}

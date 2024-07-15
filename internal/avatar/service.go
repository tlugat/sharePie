package avatar

import (
	"errors"
	"fmt"
	"sharePie-api/internal/models"
	"sharePie-api/internal/types"
)

type Service struct {
	Repository types.IAvatarRepository
}

func NewService(
	repository types.IAvatarRepository,
) types.IAvatarService {
	return &Service{
		Repository: repository,
	}
}

func (service *Service) Find() ([]models.Avatar, error) {
	avatars, err := service.Repository.Find()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to find avatars: %v", err))
	}
	return avatars, nil
}

func (service *Service) FindOne(id uint) (models.Avatar, error) {
	avatar, err := service.Repository.FindOne(id)
	if err != nil {
		return models.Avatar{}, errors.New(fmt.Sprintf("failed to find avatar with id %d: %v", id, err))
	}
	return avatar, nil
}

func (service *Service) Create(input types.CreateAvatarInput) (models.Avatar, error) {
	avatar := models.Avatar{
		Name: input.Name,
		URL:  input.URL,
	}
	newAvatar, err := service.Repository.Create(avatar)
	if err != nil {
		return models.Avatar{}, errors.New(fmt.Sprintf("failed to create avatar: %v", err))
	}
	return newAvatar, nil
}

func (service *Service) Update(id uint, input types.UpdateAvatarInput) (models.Avatar, error) {
	avatar, err := service.Repository.FindOne(id)
	if err != nil {
		return models.Avatar{}, errors.New(fmt.Sprintf("failed to find avatar with id %d: %v", id, err))
	}

	if input.Name != "" {
		avatar.Name = input.Name
	}
	if input.URL != "" {
		avatar.URL = input.URL
	}

	updatedAvatar, err := service.Repository.Update(avatar)
	if err != nil {
		return models.Avatar{}, errors.New(fmt.Sprintf("failed to update avatar with id %d: %v", id, err))
	}
	return updatedAvatar, nil
}

func (service *Service) Delete(id uint) error {
	err := service.Repository.Delete(id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete avatar with id %d: %v", id, err))
	}
	return nil
}

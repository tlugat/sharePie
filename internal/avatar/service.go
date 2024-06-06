package avatar

import (
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
		return nil, err
	}

	return avatars, nil
}

func (service *Service) FindOne(id uint) (models.Avatar, error) {
	avatar, err := service.Repository.FindOne(id)

	if err != nil {
		return models.Avatar{}, err
	}

	return avatar, nil
}

func (service *Service) Create(input types.CreateAvatarInput) (models.Avatar, error) {
	avatar := models.Avatar{
		Name: input.Name,
		URL:  input.URL,
	}

	return service.Repository.Create(avatar)
}

func (service *Service) Update(id uint, input types.UpdateAvatarInput) (models.Avatar, error) {
	avatar, err := service.Repository.FindOne(id)

	if err != nil {
		return models.Avatar{}, err
	}

	if input.Name != "" {
		avatar.Name = input.Name
	}
	if input.URL != "" {
		avatar.URL = input.URL
	}

	return service.Repository.Update(avatar)
}

func (service *Service) Delete(id uint) error {
	return service.Repository.Delete(id)
}

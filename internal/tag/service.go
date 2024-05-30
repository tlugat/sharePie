package tag

import (
	"sharePie-api/internal/models"
	"sharePie-api/internal/types"
)

type Service struct {
	Repository types.ITagRepository
}

func NewService(repository types.ITagRepository) types.ITagService {
	return &Service{Repository: repository}
}

func (service *Service) Find() ([]models.Tag, error) {
	return service.Repository.Find()
}

func (service *Service) FindOne(id uint) (models.Tag, error) {
	return service.Repository.FindOne(id)
}

func (service *Service) Create(input types.CreateTagInput) (models.Tag, error) {
	tag := models.Tag{Name: input.Name}
	return service.Repository.Create(tag)
}

func (service *Service) Update(id uint, input types.UpdateTagInput) (models.Tag, error) {
	tag, err := service.Repository.FindOne(id)

	if err != nil {
		return models.Tag{}, err
	}

	if input.Name != "" {
		tag.Name = input.Name
	}

	return service.Repository.Update(tag)
}

func (service *Service) Delete(id uint) error {
	return service.Repository.Delete(id)
}

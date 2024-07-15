package tag

import (
	"errors"
	"fmt"
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
	tags, err := service.Repository.Find()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to find tags: %v", err))
	}
	return tags, nil
}

func (service *Service) FindOne(id uint) (models.Tag, error) {
	tag, err := service.Repository.FindOne(id)
	if err != nil {
		return models.Tag{}, errors.New(fmt.Sprintf("failed to find tag with id %d: %v", id, err))
	}
	return tag, nil
}

func (service *Service) Create(input types.CreateTagInput) (models.Tag, error) {
	tag := models.Tag{Name: input.Name}
	newTag, err := service.Repository.Create(tag)
	if err != nil {
		return models.Tag{}, errors.New(fmt.Sprintf("failed to create tag: %v", err))
	}
	return newTag, nil
}

func (service *Service) Update(id uint, input types.UpdateTagInput) (models.Tag, error) {
	tag, err := service.Repository.FindOne(id)
	if err != nil {
		return models.Tag{}, errors.New(fmt.Sprintf("failed to find tag with id %d: %v", id, err))
	}

	if input.Name != "" {
		tag.Name = input.Name
	}

	updatedTag, err := service.Repository.Update(tag)
	if err != nil {
		return models.Tag{}, errors.New(fmt.Sprintf("failed to update tag with id %d: %v", id, err))
	}
	return updatedTag, nil
}

func (service *Service) Delete(id uint) error {
	err := service.Repository.Delete(id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete tag with id %d: %v", id, err))
	}
	return nil
}

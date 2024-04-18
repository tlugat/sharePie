package services

import (
	"sharePie-api/models"
	"sharePie-api/repositories"
)

type CreateTagInput struct {
	Name string `json:"name" binding:"required"`
}

type UpdateTagInput struct {
	Name string `json:"title"`
}

type ITagService interface {
	Find() ([]models.Tag, error)
	FindOne(id uint) (models.Tag, error)
	Create(input CreateTagInput) (models.Tag, error)
	Update(id uint, input UpdateTagInput) (models.Tag, error)
	Delete(id uint) error
}

type TagService struct {
	Repository repositories.ITagRepository
}

func NewTagService(repository repositories.ITagRepository) ITagService {
	return &TagService{Repository: repository}
}

func (service *TagService) Find() ([]models.Tag, error) {
	return service.Repository.Find()
}

func (service *TagService) FindOne(id uint) (models.Tag, error) {
	return service.Repository.FindOne(id)
}

func (service *TagService) Create(input CreateTagInput) (models.Tag, error) {
	tag := models.Tag{Name: input.Name}
	return service.Repository.Create(tag)
}

func (service *TagService) Update(id uint, input UpdateTagInput) (models.Tag, error) {
	tag, err := service.Repository.FindOne(id)

	if err != nil {
		return models.Tag{}, err
	}

	if input.Name != "" {
		tag.Name = input.Name
	}

	return service.Repository.Update(tag)
}

func (service *TagService) Delete(id uint) error {
	return service.Repository.Delete(id)
}

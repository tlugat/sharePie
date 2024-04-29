package category

import (
	"sharePie-api/internal/models"
)

type CreateCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCategoryInput struct {
	Name string `json:"name"`
}

type ICategoryService interface {
	Find() ([]models.Category, error)
	FindOne(id uint) (models.Category, error)
	Create(input CreateCategoryInput) (models.Category, error)
	Update(id uint, input UpdateCategoryInput) (models.Category, error)
	Delete(id uint) error
}

type Service struct {
	Repository ICategoryRepository
}

func NewService(repository ICategoryRepository) ICategoryService {
	return &Service{Repository: repository}
}

func (service *Service) Find() ([]models.Category, error) {
	return service.Repository.Find()
}

func (service *Service) FindOne(id uint) (models.Category, error) {
	return service.Repository.FindOne(id)
}

func (service *Service) Create(input CreateCategoryInput) (models.Category, error) {
	category := models.Category{Name: input.Name}
	return service.Repository.Create(category)
}

func (service *Service) Update(id uint, input UpdateCategoryInput) (models.Category, error) {
	category, err := service.Repository.FindOne(id)

	if err != nil {
		return models.Category{}, err
	}

	if input.Name != "" {
		category.Name = input.Name
	}

	return service.Repository.Update(category)
}

func (service *Service) Delete(id uint) error {
	return service.Repository.Delete(id)
}

package category

import (
	"sharePie-api/internal/models"
	"sharePie-api/internal/types"
)

type Service struct {
	Repository types.ICategoryRepository
}

func NewService(repository types.ICategoryRepository) types.ICategoryService {
	return &Service{Repository: repository}
}

func (service *Service) Find() ([]models.Category, error) {
	return service.Repository.Find()
}

func (service *Service) FindOne(id uint) (models.Category, error) {
	return service.Repository.FindOne(id)
}

func (service *Service) Create(input types.CreateCategoryInput) (models.Category, error) {
	category := models.Category{Name: input.Name}
	return service.Repository.Create(category)
}

func (service *Service) Update(id uint, input types.UpdateCategoryInput) (models.Category, error) {
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

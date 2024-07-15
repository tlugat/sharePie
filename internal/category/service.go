package category

import (
	"errors"
	"fmt"
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
	categories, err := service.Repository.Find()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to find categories: %v", err))
	}
	return categories, nil
}

func (service *Service) FindOne(id uint) (models.Category, error) {
	category, err := service.Repository.FindOne(id)
	if err != nil {
		return models.Category{}, errors.New(fmt.Sprintf("failed to find category with id %d: %v", id, err))
	}
	return category, nil
}

func (service *Service) Create(input types.CreateCategoryInput) (models.Category, error) {
	category := models.Category{Name: input.Name}
	newCategory, err := service.Repository.Create(category)
	if err != nil {
		return models.Category{}, errors.New(fmt.Sprintf("failed to create category: %v", err))
	}
	return newCategory, nil
}

func (service *Service) Update(id uint, input types.UpdateCategoryInput) (models.Category, error) {
	category, err := service.Repository.FindOne(id)
	if err != nil {
		return models.Category{}, errors.New(fmt.Sprintf("failed to find category with id %d: %v", id, err))
	}

	if input.Name != "" {
		category.Name = input.Name
	}

	updatedCategory, err := service.Repository.Update(category)
	if err != nil {
		return models.Category{}, errors.New(fmt.Sprintf("failed to update category with id %d: %v", id, err))
	}
	return updatedCategory, nil
}

func (service *Service) Delete(id uint) error {
	err := service.Repository.Delete(id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete category with id %d: %v", id, err))
	}
	return nil
}

package services

import (
	"go-project/models"
	"go-project/repositories"
)

type CreateCategoryInput struct {
	Name string `json:"title" binding:"required"`
}

type UpdateCategoryInput struct {
	Name string `json:"title"`
}

type ICategoryService interface {
	Find() ([]models.Category, error)
	FindOne(id int) (models.Category, error)
	Create(input CreateCategoryInput) (models.Category, error)
	Update(id int, input UpdateCategoryInput) (models.Category, error)
	Delete(id int) error
}

type CategoryService struct {
	Repository repositories.ICategoryRepository
}

func NewCategoryService(repository repositories.ICategoryRepository) ICategoryService {
	return &CategoryService{Repository: repository}
}

func (service *CategoryService) Find() ([]models.Category, error) {
	return service.Repository.Find()
}

func (service *CategoryService) FindOne(id int) (models.Category, error) {
	return service.Repository.FindOne(id)
}

func (service *CategoryService) Create(input CreateCategoryInput) (models.Category, error) {
	category := models.Category{Name: input.Name}
	return service.Repository.Create(category)
}

func (service *CategoryService) Update(id int, input UpdateCategoryInput) (models.Category, error) {
	category, err := service.Repository.FindOne(id)

	if err != nil {
		return models.Category{}, err
	}

	if input.Name != "" {
		category.Name = input.Name
	}

	return service.Repository.Update(category)
}

func (service *CategoryService) Delete(id int) error {
	return service.Repository.Delete(id)
}

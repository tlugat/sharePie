package services

import (
	"go-project/models"
	"go-project/repositories"
)

type CreateBookInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

type UpdateBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type IBookService interface {
	Find() ([]models.Book, error)
	FindOne(id int) (models.Book, error)
	Create(input CreateBookInput) (models.Book, error)
	Update(id int, input UpdateBookInput) (models.Book, error)
	Delete(id int) error
}

type BookService struct {
	Repository repositories.IBookRepository
}

func NewBookService(repository repositories.IBookRepository) IBookService {
	return &BookService{Repository: repository}
}

func (service *BookService) Find() ([]models.Book, error) {
	return service.Repository.Find()
}

func (service *BookService) FindOne(id int) (models.Book, error) {
	return service.Repository.FindOne(id)
}

func (service *BookService) Create(input CreateBookInput) (models.Book, error) {
	book := models.Book{Title: input.Title, Author: input.Author}
	return service.Repository.Create(book)
}

func (service *BookService) Update(id int, input UpdateBookInput) (models.Book, error) {
	book, err := service.Repository.FindOne(id)

	if err != nil {
		return models.Book{}, err
	}

	if input.Title != "" {
		book.Title = input.Title
	}
	if input.Author != "" {
		book.Author = input.Author
	}

	return service.Repository.Update(book)
}

func (service *BookService) Delete(id int) error {
	return service.Repository.Delete(id)
}

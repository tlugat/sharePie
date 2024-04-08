package repositories

import (
	"go-project/models"
	"gorm.io/gorm"
)

type IBookRepository interface {
	Find() ([]models.Book, error)
	FindOne(id int) (models.Book, error)
	Create(book models.Book) (models.Book, error)
	Update(book models.Book) (models.Book, error)
	Delete(id int) error
}

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) IBookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Find() ([]models.Book, error) {
	var books []models.Book
	result := r.db.Find(&books)
	return books, result.Error
}

func (r *BookRepository) FindOne(id int) (models.Book, error) {
	var book models.Book
	result := r.db.Where("id = ?", id).First(&book)
	return book, result.Error
}

func (r *BookRepository) Create(book models.Book) (models.Book, error) {
	result := r.db.Create(&book)
	return book, result.Error
}

func (r *BookRepository) Update(book models.Book) (models.Book, error) {
	result := r.db.Save(&book)
	return book, result.Error
}

func (r *BookRepository) Delete(id int) error {
	result := r.db.Delete(&models.Book{}, id)
	return result.Error
}

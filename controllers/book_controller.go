package controllers

import (
	"github.com/gin-gonic/gin"
	"go-project/services"
	"net/http"
	"strconv"
)

type BookController struct {
	bookService services.IBookService
}

func NewBookController(service services.IBookService) *BookController {
	return &BookController{bookService: service}
}

func (controller *BookController) FindBooks(c *gin.Context) {
	books, err := controller.bookService.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": books})
}

func (controller *BookController) FindBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	book, err := controller.bookService.FindOne(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": book})
}

func (controller *BookController) CreateBook(c *gin.Context) {
	var input services.CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	book, err := controller.bookService.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": book})
}

func (controller *BookController) UpdateBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input services.UpdateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	book, err := controller.bookService.Update(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": book})
}

func (controller *BookController) DeleteBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := controller.bookService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete book"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}

package routes

import (
	"github.com/gin-gonic/gin"
	"go-project/controllers"
	"go-project/repositories"
	"go-project/services"
	"gorm.io/gorm"
)

func InitRoutes(db *gorm.DB, route *gin.RouterGroup) {
	BookHandler(db, route)
}

func BookHandler(db *gorm.DB, route *gin.RouterGroup) {
	bookRepository := repositories.NewBookRepository(db)
	bookService := services.NewBookService(bookRepository)
	bookController := controllers.NewBookController(bookService)

	route.GET("/books", bookController.FindBooks)
	route.POST("/books", bookController.CreateBook)
	route.GET("/books/:id", bookController.FindBook)
	route.PATCH("/books/:id", bookController.UpdateBook)
	route.DELETE("/books/:id", bookController.DeleteBook)
}

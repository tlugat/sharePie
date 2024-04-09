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
	CategoryHandler(db, route)
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

func CategoryHandler(db *gorm.DB, route *gin.RouterGroup) {
	categoryRepository := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepository)
	categoryController := controllers.NewCategoryController(categoryService)

	route.GET("/categories", categoryController.FindCategories)
	route.POST("/categories", categoryController.CreateCategory)
	route.GET("/categories/:id", categoryController.FindCategory)
	route.PATCH("/categories/:id", categoryController.UpdateCategory)
	route.DELETE("/categories/:id", categoryController.DeleteCategory)
}

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
	TagHandler(db, route)
	UserHandler(db, route)
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

func TagHandler(db *gorm.DB, route *gin.RouterGroup) {
	tagRepository := repositories.NewTagRepository(db)
	tagService := services.NewTagService(tagRepository)
	tagController := controllers.NewTagController(tagService)

	route.GET("/tags", tagController.FindTags)
	route.POST("/tags", tagController.CreateTag)
	route.GET("/tags/:id", tagController.FindTag)
	route.PATCH("/tags/:id", tagController.UpdateTag)
	route.DELETE("/tags/:id", tagController.DeleteTag)
}

func UserHandler(db *gorm.DB, route *gin.RouterGroup) {
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	route.GET("/users", userController.FindUsers)
	route.GET("/users/:id", userController.FindUser)
	route.PATCH("/users/:id", userController.UpdateUser)
	route.DELETE("/users/:id", userController.DeleteUser)
}

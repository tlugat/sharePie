package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sharePie-api/controllers"
	"sharePie-api/middlewares"
	"sharePie-api/repositories"
	"sharePie-api/services"
)

func InitRoutes(db *gorm.DB, route *gin.RouterGroup) {
	authMiddleware := middlewares.RequireAuth(db)

	CategoryHandler(db, route)
	TagHandler(db, route)
	UserHandler(db, route)
	AuthHandler(db, route, authMiddleware)
	EventHandler(db, route, authMiddleware)
	ExpenseHandler(db, route, authMiddleware)
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

func AuthHandler(db *gorm.DB, route *gin.RouterGroup, middleware gin.HandlerFunc) {
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	authController := controllers.NewAuthController(userService)

	route.POST("/signup", authController.Signup)
	route.POST("/login", authController.Login)
	route.GET("/validate", middleware, authController.Validate)
}

func EventHandler(db *gorm.DB, route *gin.RouterGroup, middleware gin.HandlerFunc) {
	eventRepository := repositories.NewEventRepository(db)
	categoryRepository := repositories.NewCategoryRepository(db)
	userRepository := repositories.NewUserRepository(db)
	eventService := services.NewEventService(eventRepository, categoryRepository, userRepository)
	eventController := controllers.NewEventController(eventService)

	route.GET("/events", middleware, eventController.FindEvents)
	route.POST("/events", middleware, eventController.CreateEvent)
	route.GET("/events/:id", middleware, eventController.FindEvent)
	route.PATCH("/events/:id", middleware, eventController.UpdateEvent)
	route.DELETE("/events/:id", middleware, eventController.DeleteEvent)
}

func ExpenseHandler(db *gorm.DB, route *gin.RouterGroup, middleware gin.HandlerFunc) {
	expenseRepository := repositories.NewExpenseRepository(db)
	tagRepository := repositories.NewTagRepository(db)
	userRepository := repositories.NewUserRepository(db)
	expenseService := services.NewExpenseService(expenseRepository, tagRepository, userRepository)
	expenseController := controllers.NewExpenseController(expenseService)

	route.GET("/expenses", middleware, expenseController.FindExpenses)
	route.POST("/expenses", middleware, expenseController.CreateExpense)
	route.GET("/expenses/:id", middleware, expenseController.FindExpense)
	route.PATCH("/expenses/:id", middleware, expenseController.UpdateExpense)
	route.DELETE("/expenses/:id", middleware, expenseController.DeleteExpense)
}

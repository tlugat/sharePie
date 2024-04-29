package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sharePie-api/controllers"
	"sharePie-api/middlewares"
	"sharePie-api/repositories"
	"sharePie-api/services"
	"sharePie-api/utils"
)

func CategoryHandler(db *gorm.DB, route *gin.RouterGroup, middleware gin.HandlerFunc) {
	categoryRepository := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepository)
	categoryController := controllers.NewCategoryController(categoryService)

	route.GET("/categories", middleware, categoryController.FindCategories)
	route.POST("/categories", middleware, middlewares.AuthorizeRole(utils.AdminRole), categoryController.CreateCategory)
	route.GET("/categories/:id", middleware, categoryController.FindCategory)
	route.PATCH("/categories/:id", middleware, middlewares.AuthorizeRole(utils.AdminRole), categoryController.UpdateCategory)
	route.DELETE("/categories/:id", middleware, middlewares.AuthorizeRole(utils.AdminRole), categoryController.DeleteCategory)
}

func TagHandler(db *gorm.DB, route *gin.RouterGroup, middleware gin.HandlerFunc) {
	tagRepository := repositories.NewTagRepository(db)
	tagService := services.NewTagService(tagRepository)
	tagController := controllers.NewTagController(tagService)

	route.GET("/tags", middleware, tagController.FindTags)
	route.POST("/tags", middleware, middlewares.AuthorizeRole(utils.AdminRole), tagController.CreateTag)
	route.GET("/tags/:id", middleware, tagController.FindTag)
	route.PATCH("/tags/:id", middleware, middlewares.AuthorizeRole(utils.AdminRole), tagController.UpdateTag)
	route.DELETE("/tags/:id", middleware, middlewares.AuthorizeRole(utils.AdminRole), tagController.DeleteTag)
}

func UserHandler(db *gorm.DB, route *gin.RouterGroup, middleware gin.HandlerFunc) {
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	route.GET("/users", middleware, middlewares.AuthorizeRole(utils.AdminRole), userController.FindUsers)
	route.GET("/users/:id", middleware, userController.FindUser)
	route.PATCH("/users/:id", middleware, userController.UpdateUser)
	route.DELETE("/users/:id", middleware, middlewares.AuthorizeRole(utils.AdminRole), userController.DeleteUser)
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
	expenseRepository := repositories.NewExpenseRepository(db)
	eventService := services.NewEventService(eventRepository, categoryRepository, userRepository, expenseRepository)
	eventBalanceService := services.NewEventBalanceService(eventRepository, expenseRepository, userRepository)
	eventController := controllers.NewEventController(eventService, eventBalanceService)

	route.GET("/events", middleware, eventController.FindEvents)
	route.POST("/events", middleware, eventController.CreateEvent)
	route.GET("/events/:id", middleware, eventController.FindEvent)
	route.PATCH("/events/:id", middleware, eventController.UpdateEvent)
	route.DELETE("/events/:id", middleware, eventController.DeleteEvent)
	route.GET("/events/:id/summary", middleware, eventController.GetEventBalanceSummary)
	route.GET("/events/:id/users", middleware, eventController.GetEventUsers)
	route.POST("/events/join", middleware, eventController.JoinEvent)
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

func AchievementHandler(db *gorm.DB, route *gin.RouterGroup, middleware gin.HandlerFunc) {
	achievementRepository := repositories.NewAchievementRepository(db)
	achievementService := services.NewAchievementService(achievementRepository)
	achievementController := controllers.NewAchievementController(achievementService)

	route.GET("/achievements", middleware, achievementController.FindAchievements)
	route.POST("/achievements", middleware, achievementController.CreateAchievement)
	route.GET("/achievements/:id", middleware, achievementController.FindAchievement)
	route.PATCH("/achievements/:id", middleware, achievementController.UpdateAchievement)
	route.DELETE("/achievements/:id", middleware, achievementController.DeleteAchievement)
}

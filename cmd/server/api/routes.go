package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sharePie-api/internal/achievement"
	"sharePie-api/internal/auth"
	"sharePie-api/internal/auth/middleware"
	"sharePie-api/internal/category"
	"sharePie-api/internal/event"
	"sharePie-api/internal/expense"
	"sharePie-api/internal/tag"
	"sharePie-api/internal/user"
	"sharePie-api/pkg/constants"
)

func CategoryHandler(db *gorm.DB, route *gin.RouterGroup) {
	categoryRepository := category.NewRepository(db)
	categoryService := category.NewService(categoryRepository)
	categoryController := category.NewController(categoryService)

	route.GET("/categories", middleware.IsAuthenticated(db), categoryController.FindCategories)
	route.POST("/categories", middleware.IsAuthenticated(db), middleware.IsGranted(constants.AdminRole), categoryController.CreateCategory)
	route.GET("/categories/:id", middleware.IsAuthenticated(db), categoryController.FindCategory)
	route.PATCH("/categories/:id", middleware.IsAuthenticated(db), middleware.IsGranted(constants.AdminRole), categoryController.UpdateCategory)
	route.DELETE("/categories/:id", middleware.IsAuthenticated(db), middleware.IsGranted(constants.AdminRole), categoryController.DeleteCategory)
}

func TagHandler(db *gorm.DB, route *gin.RouterGroup) {
	tagRepository := tag.NewRepository(db)
	tagService := tag.NewService(tagRepository)
	tagController := tag.NewController(tagService)

	route.GET("/tags", middleware.IsAuthenticated(db), tagController.FindTags)
	route.POST("/tags", middleware.IsAuthenticated(db), middleware.IsGranted(constants.AdminRole), tagController.CreateTag)
	route.GET("/tags/:id", middleware.IsAuthenticated(db), tagController.FindTag)
	route.PATCH("/tags/:id", middleware.IsAuthenticated(db), middleware.IsGranted(constants.AdminRole), tagController.UpdateTag)
	route.DELETE("/tags/:id", middleware.IsAuthenticated(db), middleware.IsGranted(constants.AdminRole), tagController.DeleteTag)
}

func UserHandler(db *gorm.DB, route *gin.RouterGroup) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userController := user.NewController(userService)

	route.GET("/users", middleware.IsAuthenticated(db), middleware.IsGranted(constants.AdminRole), userController.FindUsers)
	route.GET("/users/:id", middleware.IsAuthenticated(db), userController.FindUser)
	route.PATCH("/users/:id", middleware.IsAuthenticated(db), userController.UpdateUser)
	route.DELETE("/users/:id", middleware.IsAuthenticated(db), middleware.IsGranted(constants.AdminRole), userController.DeleteUser)
}

func AuthHandler(db *gorm.DB, route *gin.RouterGroup) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authController := auth.NewController(userService)

	route.POST("/signup", authController.Signup)
	route.POST("/login", authController.Login)
	route.GET("/validate", middleware.IsAuthenticated(db), authController.Validate)
}

func EventHandler(db *gorm.DB, route *gin.RouterGroup) {
	eventRepository := event.NewRepository(db)
	categoryRepository := category.NewRepository(db)
	userRepository := user.NewRepository(db)
	expenseRepository := expense.NewRepository(db)
	eventService := event.NewService(eventRepository, categoryRepository, userRepository, expenseRepository)
	eventBalanceService := event.NewBalanceService(eventRepository, expenseRepository, userRepository)
	eventController := event.NewController(eventService, eventBalanceService)

	route.GET("/events", middleware.IsAuthenticated(db), eventController.FindEvents)
	route.POST("/events", middleware.IsAuthenticated(db), eventController.CreateEvent)
	route.GET("/events/:id", middleware.IsAuthenticated(db), eventController.FindEvent)
	route.PATCH("/events/:id", middleware.IsAuthenticated(db), eventController.UpdateEvent)
	route.DELETE("/events/:id", middleware.IsAuthenticated(db), eventController.DeleteEvent)
	route.GET("/events/:id/summary", middleware.IsAuthenticated(db), eventController.GetEventBalanceSummary)
	route.GET("/events/:id/users", middleware.IsAuthenticated(db), eventController.GetEventUsers)
	route.POST("/events/join", middleware.IsAuthenticated(db), eventController.JoinEvent)
}

func ExpenseHandler(db *gorm.DB, route *gin.RouterGroup) {
	expenseRepository := expense.NewRepository(db)
	tagRepository := tag.NewRepository(db)
	userRepository := user.NewRepository(db)
	expenseService := expense.NewService(expenseRepository, tagRepository, userRepository)
	expenseController := expense.NewController(expenseService)

	route.GET("/expenses", middleware.IsAuthenticated(db), expenseController.FindExpenses)
	route.POST("/expenses", middleware.IsAuthenticated(db), expenseController.CreateExpense)
	route.GET("/expenses/:id", middleware.IsAuthenticated(db), expenseController.FindExpense)
	route.PATCH("/expenses/:id", middleware.IsAuthenticated(db), expenseController.UpdateExpense)
	route.DELETE("/expenses/:id", middleware.IsAuthenticated(db), expenseController.DeleteExpense)
}

func AchievementHandler(db *gorm.DB, route *gin.RouterGroup) {
	achievementRepository := achievement.NewRepository(db)
	achievementService := achievement.NewService(achievementRepository)
	achievementController := achievement.NewController(achievementService)

	route.GET("/achievements", middleware.IsAuthenticated(db), achievementController.FindAchievements)
	route.POST("/achievements", middleware.IsAuthenticated(db), achievementController.CreateAchievement)
	route.GET("/achievements/:id", middleware.IsAuthenticated(db), achievementController.FindAchievement)
	route.PATCH("/achievements/:id", middleware.IsAuthenticated(db), achievementController.UpdateAchievement)
	route.DELETE("/achievements/:id", middleware.IsAuthenticated(db), achievementController.DeleteAchievement)
}

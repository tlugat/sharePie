package api

import (
	"sharePie-api/internal/achievement"
	"sharePie-api/internal/auth"
	"sharePie-api/internal/auth/middleware"
	"sharePie-api/internal/avatar"
	"sharePie-api/internal/category"
	"sharePie-api/internal/event"
	"sharePie-api/internal/expense"
	"sharePie-api/internal/participant"
	"sharePie-api/internal/payer"
	"sharePie-api/internal/refund"
	"sharePie-api/internal/tag"
	"sharePie-api/internal/user"
	"sharePie-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CategoryHandler(db *gorm.DB, route *gin.RouterGroup) {
	categoryRepository := category.NewRepository(db)
	categoryService := category.NewService(categoryRepository)
	categoryController := category.NewController(categoryService)

	route.GET("/categories", middleware.IsAuthenticated(db), categoryController.FindCategories)
	route.POST("/categories", middleware.IsAuthenticated(db), middleware.IsGranted(utils.AdminRole), categoryController.CreateCategory)
	route.GET("/categories/:id", middleware.IsAuthenticated(db), categoryController.FindCategory)
	route.PATCH("/categories/:id", middleware.IsAuthenticated(db), middleware.IsGranted(utils.AdminRole), categoryController.UpdateCategory)
	route.DELETE("/categories/:id", middleware.IsAuthenticated(db), middleware.IsGranted(utils.AdminRole), categoryController.DeleteCategory)
}

func TagHandler(db *gorm.DB, route *gin.RouterGroup) {
	tagRepository := tag.NewRepository(db)
	tagService := tag.NewService(tagRepository)
	tagController := tag.NewController(tagService)

	route.GET("/tags", middleware.IsAuthenticated(db), tagController.FindTags)
	route.POST("/tags", middleware.IsAuthenticated(db), middleware.IsGranted(utils.AdminRole), tagController.CreateTag)
	route.GET("/tags/:id", middleware.IsAuthenticated(db), tagController.FindTag)
	route.PATCH("/tags/:id", middleware.IsAuthenticated(db), middleware.IsGranted(utils.AdminRole), tagController.UpdateTag)
	route.DELETE("/tags/:id", middleware.IsAuthenticated(db), middleware.IsGranted(utils.AdminRole), tagController.DeleteTag)
}

func UserHandler(db *gorm.DB, route *gin.RouterGroup) {
	userRepository := user.NewRepository(db)
	avatarRepository := avatar.NewRepository(db)
	userService := user.NewService(userRepository, avatarRepository)
	userController := user.NewController(userService)

	route.POST("/users", middleware.IsAuthenticated(db), middleware.IsGranted(utils.AdminRole), userController.CreateUser)
	route.GET("/users", middleware.IsAuthenticated(db), middleware.IsGranted(utils.AdminRole), userController.FindUsers)
	route.GET("/users/:id", middleware.IsAuthenticated(db), userController.FindUser)
	route.PATCH("/users/me", middleware.IsAuthenticated(db), userController.UpdateCurrentUser)
	route.PATCH("/users/firebase_token", middleware.IsAuthenticated(db), userController.UpdateCurrentUserFirebaseToken)
	route.PATCH("/users/:id", middleware.IsAuthenticated(db), middleware.IsGranted(utils.AdminRole), userController.UpdateUser)
	route.DELETE("/users/:id", middleware.IsAuthenticated(db), middleware.IsGranted(utils.AdminRole), userController.DeleteUser)
	route.GET("/users/me", middleware.IsAuthenticated(db), userController.GetUserFromToken)
}

func AuthHandler(db *gorm.DB, route *gin.RouterGroup) {
	userRepository := user.NewRepository(db)
	avatarRepository := avatar.NewRepository(db)
	userService := user.NewService(userRepository, avatarRepository)
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
	refundRepository := refund.NewRepository(db)
	eventService := event.NewService(eventRepository, categoryRepository, userRepository, expenseRepository, refundRepository)
	eventController := event.NewController(eventService)

	route.GET("/events", middleware.IsAuthenticated(db), eventController.FindEvents)
	route.POST("/events", middleware.IsAuthenticated(db), eventController.CreateEvent)
	route.GET("/events/:id", middleware.IsAuthenticated(db), middleware.EventIsUserPartOfEvent(db), eventController.FindEvent)
	route.PATCH("/events/:id", middleware.IsAuthenticated(db), middleware.IsEventActive(db), middleware.IsEventAuthor(db), eventController.UpdateEvent)
	route.PATCH("/events/:id/state", middleware.IsAuthenticated(db), middleware.IsEventAuthor(db), eventController.UpdateEventState)
	route.DELETE("/events/:id", middleware.IsAuthenticated(db), eventController.DeleteEvent)
	route.GET("/events/:id/balances", middleware.IsAuthenticated(db), eventController.GetEventBalances)
	route.GET("/events/:id/transactions", middleware.IsAuthenticated(db), eventController.GetEventTransactions)
	route.GET("/events/:id/users", middleware.IsAuthenticated(db), eventController.GetEventUsers)
	route.POST("/events/join", middleware.IsAuthenticated(db), eventController.JoinEvent)
	route.GET("/events/:id/expenses", middleware.IsAuthenticated(db), eventController.GetExpenses)
}

func ExpenseHandler(db *gorm.DB, route *gin.RouterGroup) {
	expenseRepository := expense.NewRepository(db)
	tagRepository := tag.NewRepository(db)
	userRepository := user.NewRepository(db)
	participantRepository := participant.NewRepository(db)
	payerRepository := payer.NewRepository(db)
	eventRepository := event.NewRepository(db)
	refundRepository := refund.NewRepository(db)
	eventService := event.NewService(eventRepository, nil, userRepository, nil, refundRepository)
	expenseService := expense.NewService(
		expenseRepository,
		tagRepository,
		userRepository,
		participantRepository,
		payerRepository,
		eventService,
	)
	expenseController := expense.NewController(expenseService)

	route.GET("/expenses", middleware.IsAuthenticated(db), middleware.ExpenseIsUserPartOfEvent(db), expenseController.FindExpenses)
	route.POST("/expenses", middleware.IsAuthenticated(db), expenseController.CreateExpense)
	route.GET("/expenses/:id", middleware.IsAuthenticated(db), middleware.ExpenseIsUserPartOfEvent(db), expenseController.FindExpense)
	route.PATCH("/expenses/:id", middleware.IsAuthenticated(db), middleware.ExpenseIsUserPartOfEvent(db), expenseController.UpdateExpense)
	route.DELETE("/expenses/:id", middleware.IsAuthenticated(db), middleware.ExpenseIsUserPartOfEvent(db), expenseController.DeleteExpense)
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

func AvatarHandler(db *gorm.DB, route *gin.RouterGroup) {
	avatarRepository := avatar.NewRepository(db)
	avatarService := avatar.NewService(avatarRepository)
	avatarController := avatar.NewController(avatarService)

	route.GET("/avatars", middleware.IsAuthenticated(db), avatarController.FindAvatars)
	route.POST("/avatars", middleware.IsAuthenticated(db), middleware.IsGranted(utils.AdminRole), avatarController.CreateAvatar)
	route.GET("/avatars/:id", middleware.IsAuthenticated(db), avatarController.FindAvatar)
	route.PATCH("/avatars/:id", middleware.IsAuthenticated(db), middleware.IsGranted(utils.AdminRole), avatarController.UpdateAvatar)
	route.DELETE("/avatars/:id", middleware.IsAuthenticated(db), middleware.IsGranted(utils.AdminRole), avatarController.DeleteAvatar)
}

func RefundHandler(db *gorm.DB, route *gin.RouterGroup) {
	refundRepository := refund.NewRepository(db)
	userRepository := user.NewRepository(db)
	eventRepository := event.NewRepository(db)
	eventService := event.NewService(eventRepository, nil, userRepository, nil, refundRepository)
	refundService := refund.NewService(refundRepository, userRepository, eventService)
	refundController := refund.NewController(refundService)

	route.GET("/refunds", middleware.IsAuthenticated(db), refundController.FindRefunds)
	route.GET("/refunds/:id", middleware.IsAuthenticated(db), middleware.RefundIsUserPartOfEvent(db), refundController.FindRefund)
	route.DELETE("/refunds/:id", middleware.IsAuthenticated(db), middleware.RefundIsUserPartOfEvent(db), refundController.DeleteRefund)
}

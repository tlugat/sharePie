package main

import (
	"context"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	api2 "sharePie-api/cmd/server/api"
	_ "sharePie-api/docs"
	"sharePie-api/internal/auth/middleware"
	"sharePie-api/pkg/config/database"
	"sharePie-api/pkg/config/env"
	"sharePie-api/pkg/config/thirdparty/cloudinary"
	firebase2 "sharePie-api/pkg/config/thirdparty/firebase"
	websocket2 "sharePie-api/pkg/config/thirdparty/websocket"
	"syscall"
	"time"
)

// @title SharePie API
// @version 1.0
// @description This is the API of SharePie app. You can visit the GitHub repository at https://github.com/tlugat/sharePie-api

// @host localhost:8080
// @SecurityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @BasePath /api/v1
// @query.collection.format multi
func main() {
	envMode := os.Getenv("ENV")
	if envMode != "production" {
		err := env.Load()
		if err != nil {
			log.Fatalf("Failed to load environment variables : %v", err)
		}
	}

	err := firebase2.InitFirebase()
	if err != nil {
		log.Fatalf("Failed to initialize Firebase : %v", err)
	}

	db, err := database.NewPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to database : %v", err)
	}

	err = cloudinary.NewCloudinaryClient()
	if err != nil {
		log.Fatalf("Failed to connect to Cloudinary : %v", err)
	}

	r := gin.Default()

	// Configure CORS
	corsConfig := cors.Config{
		AllowOrigins:     []string{os.Getenv("CORS_ORIGIN")},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))
	api := r.Group("api/v1")

	hub := websocket2.NewHub(db)
	go hub.Run()

	api.GET("/ws/events/:eventId", middleware.WSIsAuthenticated(db), middleware.IsEventActive(db), func(c *gin.Context) {
		websocket2.ServeWs(hub, c)
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	api2.CategoryHandler(db, api)
	api2.TagHandler(db, api)
	api2.UserHandler(db, api)
	api2.AuthHandler(db, api)
	api2.EventHandler(db, api)
	api2.ExpenseHandler(db, api)
	api2.AchievementHandler(db, api)
	api2.AvatarHandler(db, api)
	api2.RefundHandler(db, api)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		} else if err == nil {
			log.Println("Server started on port 8080")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can"t be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

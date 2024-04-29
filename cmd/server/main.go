package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	"sharePie-api/pkg/config/thirdparty"
	"sharePie-api/pkg/config/thirdparty/cloudinary"
	"syscall"
	"time"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// @title SharePie API
// @version 1.0
// @description This is the API of SharePie app. You can visit the GitHub repository at https://github.com/tlugat/sharePie-api

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi
func main() {
	err := env.Load()
	if err != nil {
		log.Fatalf("Failed to load environment variables : %v", err)
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

	api := r.Group("api/v1")

	api.GET("/ws", middleware.IsAuthenticated(db), func(c *gin.Context) {
		conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to connect to websocket server"})
			return
		}
		defer conn.Close()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("error: %v", err)
				}
				break
			}

			var event thirdparty.Event
			if err := json.Unmarshal(msg, &event); err != nil {
				fmt.Println("Error unmarshalling event:", err)
				continue
			}

			thirdparty.HandleWebsocketEvent(conn, event, db)
		}
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	api2.CategoryHandler(db, api)
	api2.TagHandler(db, api)
	api2.UserHandler(db, api)
	api2.AuthHandler(db, api)
	api2.EventHandler(db, api)
	api2.ExpenseHandler(db, api)
	api2.AchievementHandler(db, api)

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
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
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

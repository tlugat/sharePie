package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"sharePie-api/configs"
	_ "sharePie-api/docs"
	"sharePie-api/middlewares"
	"sharePie-api/routes"
	"sharePie-api/ws"
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
	configs.LoadEnv()
	db := configs.ConnectDB()

	r := gin.Default()

	authMiddleware := middlewares.RequireAuth(db)

	api := r.Group("api/v1")

	api.GET("/ws", func(c *gin.Context) {
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

			var event ws.Event
			if err := json.Unmarshal(msg, &event); err != nil {
				fmt.Println("Error unmarshalling event:", err)
				continue
			}

			ws.HandleEvent(conn, event)
		}
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	routes.CategoryHandler(db, api)
	routes.TagHandler(db, api)
	routes.UserHandler(db, api)
	routes.AuthHandler(db, api, authMiddleware)
	routes.EventHandler(db, api, authMiddleware)
	routes.ExpenseHandler(db, api, authMiddleware)

	err := r.Run()

	if err != nil {
		return
	}
}

package main

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"sharePie-api/configs"
	_ "sharePie-api/docs"
	"sharePie-api/routes"
)

func main() {
	configs.LoadEnv()
	DB := configs.ConnectDB()

	initRouter(DB)
}

// @title SharePie API
// @version 1.0
// @description This is the API of SharePie app. You can visit the GitHub repository at https://github.com/tlugat/sharePie-api

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi
func initRouter(db *gorm.DB) {
	r := gin.Default()

	api := r.Group("api/v1")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	routes.InitRoutes(db, api)

	err := r.Run()

	if err != nil {
		return
	}
}

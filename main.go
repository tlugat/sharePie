package main

import (
	"github.com/gin-gonic/gin"
	"go-project/configs"
	"go-project/routes"
	"gorm.io/gorm"
)

func main() {
	configs.LoadEnv()
	DB := configs.ConnectDB()

	initRouter(DB)
}

func initRouter(db *gorm.DB) {
	r := gin.Default()

	api := r.Group("api/v1")

	routes.InitRoutes(db, api)

	err := r.Run()

	if err != nil {
		return
	}
}

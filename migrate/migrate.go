package main

import (
	"go-project/configs"
	"go-project/models"
)

func main() {
	configs.LoadEnv()
	db := configs.ConnectDB()
	db.AutoMigrate(&models.User{})
}

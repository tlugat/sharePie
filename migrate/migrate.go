package main

import (
	"go-project/configs"
	"go-project/models"
)

func main() {
	configs.LoadEnv()
	db := configs.ConnectDB()
	err := db.AutoMigrate(
		&models.User{},
		&models.Tag{},
		&models.Category{},
		&models.Expense{},
		&models.Tag{},
		&models.Event{},
	)
	if err != nil {
		panic(err)
		return
	}
}

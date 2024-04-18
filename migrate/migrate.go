package main

import (
	"sharePie-api/configs"
	"sharePie-api/models"
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

package main

import (
	"log"
	"sharePie-api/configs"
	"sharePie-api/models"
)

func main() {
	configs.LoadEnv()
	db, err := configs.ConnectDB()

	if err != nil {
		log.Fatalf("Failed to connect to database : %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Tag{},
		&models.Category{},
		&models.Expense{},
		&models.Tag{},
		&models.Event{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database : %v", err)
	}
}

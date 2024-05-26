package main

import (
	"log"
	"sharePie-api/internal/models"
	"sharePie-api/pkg/config/database"
	"sharePie-api/pkg/config/env"
)

func main() {
	err := env.Load()
	if err != nil {
		log.Fatalf("Failed to load environment variables : %v", err)
	}

	db, err := database.NewPostgres()
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
		&models.Participant{},
		&models.Payer{},
		&models.Achievement{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database : %v", err)
	}
}

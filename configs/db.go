package configs

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func ConnectDB() (*gorm.DB, error) {
	dbUrl := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	return db, err
}

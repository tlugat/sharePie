package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"sharePie-api/internal/models"
	"sharePie-api/pkg/config/database"
	"sharePie-api/pkg/config/env"

	"gorm.io/gorm"
)

var tableModels = map[string]interface{}{
	"users":             &models.User{},
	"tags":              &models.Tag{},
	"categories":        &models.Category{},
	"expenses":          &models.Expense{},
	"events":            &models.Event{},
	"participants":      &models.Participant{},
	"payers":            &models.Payer{},
	"achievements":      &models.Achievement{},
	"expense_users":     "expense_users",
	"event_users":       "event_users",
	"user_achievements": "user_achievements",
}

func main() {
	err := env.Load()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	db, err := database.NewPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the name of the table to drop: ")
	tableName, _ := reader.ReadString('\n')
	tableName = strings.TrimSpace(tableName)

	model, exists := tableModels[tableName]
	if !exists {
		log.Fatalf("No model found for table name: %s", tableName)
	}

	dropTable(db, model)
}

func dropTable(db *gorm.DB, model interface{}) {
	err := db.Migrator().DropTable(model)
	if err != nil {
		log.Fatalf("Error dropping table: %v", err)
	}
	fmt.Println("Table dropped successfully")
}

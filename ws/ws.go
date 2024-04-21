package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"sharePie-api/repositories"
	"sharePie-api/services"
)

type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func HandleEvent(conn *websocket.Conn, event Event, db *gorm.DB) {
	eventRepository := repositories.NewEventRepository(db)
	expenseRepository := repositories.NewExpenseRepository(db)
	userRepository := repositories.NewUserRepository(db)
	eventBalanceService := services.NewEventBalanceService(eventRepository, expenseRepository, userRepository)

	switch event.Type {
	case "join":
		handleUserJoined(conn, eventBalanceService)
	case "leave":
		handleUserLeft(conn)
	default:
		fmt.Println("Unhandled event type:", event.Type)
	}
}

func handleUserJoined(conn *websocket.Conn, eventBalanceService services.IEventBalanceService) {
	fmt.Println("User joined")
}

func handleUserLeft(conn *websocket.Conn) {
	fmt.Println("User left")
}

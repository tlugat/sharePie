package thirdparty

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"sharePie-api/internal/event"
	"sharePie-api/internal/expense"
	"sharePie-api/internal/user"
)

type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func HandleWebsocketEvent(conn *websocket.Conn, evt Event, db *gorm.DB) {
	eventRepository := event.NewRepository(db)
	expenseRepository := expense.NewRepository(db)
	userRepository := user.NewRepository(db)
	eventBalanceService := event.NewBalanceService(eventRepository, expenseRepository, userRepository)

	switch evt.Type {
	case "join":
		handleUserJoined(conn, eventBalanceService)
	case "leave":
		handleUserLeft(conn)
	default:
		fmt.Println("Unhandled event type:", evt.Type)
	}
}

func handleUserJoined(conn *websocket.Conn, eventBalanceService event.IEventBalanceService) {
	fmt.Println("User joined")
}

func handleUserLeft(conn *websocket.Conn) {
	fmt.Println("User left")
}

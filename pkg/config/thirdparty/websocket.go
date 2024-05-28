package thirdparty

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func HandleWebsocketEvent(conn *websocket.Conn, event Event, db *gorm.DB) {

	switch event.Type {
	case "join":
		handleUserJoined(conn)
	case "leave":
		handleUserLeft(conn)
	default:
		fmt.Println("Unhandled event type:", event.Type)
	}
}

func handleUserJoined(conn *websocket.Conn) {
	fmt.Println("User joined")
}

func handleUserLeft(conn *websocket.Conn) {
	fmt.Println("User left")
}

package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func HandleEvent(conn *websocket.Conn, event Event) {
	switch event.Type {
	case "message":
		handleMessage(conn, event.Data)
	case "join":
		handleUserJoined(conn)
	case "leave":
		handleUserLeft(conn)
	default:
		fmt.Println("Unhandled event type:", event.Type)
	}
}

func handleMessage(conn *websocket.Conn, data json.RawMessage) {
	fmt.Println("Message received:", string(data))
}

func handleUserJoined(conn *websocket.Conn) {
	fmt.Println("User joined")
}

func handleUserLeft(conn *websocket.Conn) {
	fmt.Println("User left")
}

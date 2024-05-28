package thirdparty

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/url"
	"sharePie-api/internal/category"
	"sharePie-api/internal/event"
	"sharePie-api/internal/expense"
	"sharePie-api/internal/types"
	"sharePie-api/internal/user"
	"strconv"
)

type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func HandleWebsocketEvent(conn *websocket.Conn, evt Event, db *gorm.DB, queryParams url.Values) {
	eventRepository := event.NewRepository(db)
	categoryRepository := category.NewRepository(db)
	userRepository := user.NewRepository(db)
	expenseRepository := expense.NewRepository(db)
	eventService := event.NewService(eventRepository, categoryRepository, userRepository, expenseRepository)

	switch evt.Type {
	case "updateEvent":
		handleUpdateEvent(conn, eventService, evt.Data, queryParams)
	default:
		fmt.Println("Unhandled event type:", evt.Type)
	}
}

func handleUpdateEvent(conn *websocket.Conn, eventService types.IEventService, data json.RawMessage, queryParams url.Values) {
	var input types.UpdateEventInput
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Println("Failed to unmarshal data:", err)
		return
	}

	eventIdStr := queryParams.Get("eventId")
	if eventIdStr == "" {
		fmt.Println("eventId is missing in query parameters")
		return
	}

	eventId, err := strconv.ParseUint(eventIdStr, 10, 32)
	if err != nil {
		fmt.Println("Invalid eventId:", err)
		return
	}

	updatedEvent, err := eventService.Update(uint(eventId), input)
	if err != nil {
		fmt.Println("Failed to update event:", err)
		return
	}

	eventJson, err := json.Marshal(updatedEvent)
	if err != nil {
		fmt.Println("Error marshalling event:", err)
		return
	}

	response := Event{
		Type: "event",
		Data: eventJson,
	}
	if err := conn.WriteJSON(response); err != nil {
		fmt.Println("Failed to write JSON response:", err)
		return
	}
}

package websocket

import (
	"gorm.io/gorm"
	"sharePie-api/internal/category"
	"sharePie-api/internal/event"
	"sharePie-api/internal/expense"
	"sharePie-api/internal/types"
	"sharePie-api/internal/user"
)

type Hub struct {
	rooms        map[string]*Room
	register     chan *Client
	unregister   chan *Client
	eventService types.IEventService
}

func NewHub(db *gorm.DB) *Hub {
	eventRepository := event.NewRepository(db)
	categoryRepository := category.NewRepository(db)
	userRepository := user.NewRepository(db)
	expenseRepository := expense.NewRepository(db)
	eventService := event.NewService(eventRepository, categoryRepository, userRepository, expenseRepository)

	return &Hub{
		rooms:        make(map[string]*Room),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		eventService: eventService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			room, ok := h.rooms[client.room]
			if !ok {
				room = NewRoom(client.room)
				h.rooms[client.room] = room
				go room.Run()
			}
			room.register <- client
		case client := <-h.unregister:
			if room, ok := h.rooms[client.room]; ok {
				room.unregister <- client
			}
		}
	}
}

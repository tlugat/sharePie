package websocket

import (
	"gorm.io/gorm"
	"sharePie-api/internal/category"
	"sharePie-api/internal/event"
	"sharePie-api/internal/expense"
	"sharePie-api/internal/participant"
	"sharePie-api/internal/payer"
	"sharePie-api/internal/refund"
	"sharePie-api/internal/tag"
	"sharePie-api/internal/types"
	"sharePie-api/internal/user"
)

type Hub struct {
	rooms          map[string]*Room
	register       chan *Client
	unregister     chan *Client
	eventService   types.IEventService
	expenseService types.IExpenseService
	refundService  types.IRefundService
}

func NewHub(db *gorm.DB) *Hub {
	eventRepository := event.NewRepository(db)
	categoryRepository := category.NewRepository(db)
	userRepository := user.NewRepository(db)
	expenseRepository := expense.NewRepository(db)
	tagRepository := tag.NewRepository(db)
	participantRepository := participant.NewRepository(db)
	payerRepository := payer.NewRepository(db)
	refundRepository := refund.NewRepository(db)
	eventService := event.NewService(eventRepository, categoryRepository, userRepository, expenseRepository, refundRepository)
	refundService := refund.NewService(refundRepository, userRepository, eventService)
	expenseService := expense.NewService(
		expenseRepository,
		tagRepository,
		userRepository,
		participantRepository,
		payerRepository,
		eventService,
	)

	return &Hub{
		rooms:          make(map[string]*Room),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		eventService:   eventService,
		expenseService: expenseService,
		refundService:  refundService,
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

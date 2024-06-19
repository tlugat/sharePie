package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sharePie-api/internal/auth"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs(hub *Hub, c *gin.Context) {
	room := c.Param("eventId")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return
	}
	client := &Client{
		hub:  hub,
		conn: conn,
		room: room,
		send: make(chan []byte, 256),
		user: user,
	}
	client.hub.register <- client

	initRoomMessages(client)

	go client.writePump()
	go client.readPump()
}

// Fetch event, users, and expenses of the event and send the 'event', 'users', and 'expenses' events to the room
func initRoomMessages(client *Client) {
	eventId, err := strconv.ParseUint(client.room, 10, 32)
	if err == nil {
		// Fetch and send event
		event, err := client.hub.eventService.FindOne(uint(eventId))
		if err == nil {
			eventJSON, err := json.Marshal(event)
			if err == nil {
				eventMessage, err := json.Marshal(Message{
					Type:    "event",
					Payload: eventJSON,
				})
				if err == nil {
					client.hub.rooms[client.room].broadcast <- eventMessage
				} else {
					fmt.Println("Failed to marshal event message:", err)
				}
			} else {
				fmt.Println("Failed to marshal event:", err)
			}
		} else {
			fmt.Println("Failed to get event:", err)
		}

		// Fetch and send users
		users, err := client.hub.eventService.GetUsers(uint(eventId))
		if err == nil {
			usersJSON, err := json.Marshal(users)
			if err == nil {
				usersMessage, err := json.Marshal(Message{
					Type:    "users",
					Payload: usersJSON,
				})
				if err == nil {
					client.hub.rooms[client.room].broadcast <- usersMessage
				} else {
					fmt.Println("Failed to marshal users message:", err)
				}
			} else {
				fmt.Println("Failed to marshal users:", err)
			}
		} else {
			fmt.Println("Failed to get event users:", err)
		}

		// Fetch and send expenses
		expenses, err := client.hub.expenseService.FindByEventId(uint(eventId))
		if err == nil {
			expensesJSON, err := json.Marshal(expenses)
			if err == nil {
				expensesMessage, err := json.Marshal(Message{
					Type:    "expenses",
					Payload: expensesJSON,
				})
				if err == nil {
					client.hub.rooms[client.room].broadcast <- expensesMessage
				} else {
					fmt.Println("Failed to marshal expenses message:", err)
				}
			} else {
				fmt.Println("Failed to marshal expenses:", err)
			}
		} else {
			fmt.Println("Failed to get event expenses:", err)
		}

		balances, err := client.hub.eventService.GetBalances(event)
		if err == nil {
			balancesJSON, err := json.Marshal(balances)
			if err == nil {
				balancesMessage, err := json.Marshal(Message{
					Type:    "balances",
					Payload: balancesJSON,
				})
				if err == nil {
					client.hub.rooms[client.room].broadcast <- balancesMessage
				} else {
					fmt.Println("Failed to marshal balances message:", err)
				}
			} else {
				fmt.Println("Failed to marshal balances:", err)
			}
		} else {
			fmt.Println("Failed to get event balances:", err)
		}
	}
}

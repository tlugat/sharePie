package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"sharePie-api/internal/auth/middleware"
	"sharePie-api/internal/models"
	"sharePie-api/internal/types"

	"sharePie-api/pkg/config/thirdparty/firebase"
	"strconv"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	room string
	send chan []byte
	user models.User
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("error unmarshaling message:", err)
			return
		}
		c.handleMessage(msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, message)
		case <-ticker.C:
			c.conn.WriteMessage(websocket.PingMessage, []byte{})
		}
	}
}

func (c *Client) handleMessage(message Message) {
	switch message.Type {
	case "updateEvent":
		c.updateEvent(message.Payload)
	case "createExpense":
		c.createExpense(message.Payload)
	case "updateExpense":
		c.updateExpense(message.Payload)
	case "deleteExpense":
		c.deleteExpense(message.Payload)
	case "createRefund":
		c.createRefund(message.Payload)
	case "updateRefund":
		c.updateRefund(message.Payload)
	case "deleteRefund":
		c.deleteRefund(message.Payload)
	default:
		log.Println("unknown event type:", message.Type)
	}
}

func (c *Client) updateEvent(payload json.RawMessage) {
	var input types.UpdateEventInput
	if err := json.Unmarshal(payload, &input); err != nil {
		fmt.Println("Failed to unmarshal data:", err)
		return
	}
	eventId, err := strconv.ParseUint(c.room, 10, 32)
	if err != nil {
		fmt.Println("Invalid eventId:", err)
		return
	}

	event, err := c.hub.eventService.FindOne(uint(eventId))
	if err != nil {
		fmt.Println("Failed to get event:", err)
		return
	}

	if !middleware.IsUserEventAuthor(c.user, event) {
		fmt.Println("User is not the author of the event")
		return
	}

	_, err = c.hub.eventService.Update(uint(eventId), input)
	if err != nil {
		fmt.Println("Failed to update event:", err)
		return
	}

	c.refreshEvent()
	c.refreshUsers(event.ID)
}

func (c *Client) createExpense(payload json.RawMessage) {
	var input types.CreateExpenseInput
	if err := json.Unmarshal(payload, &input); err != nil {
		fmt.Println("Failed to unmarshal data:", err)
		return
	}

	expense, err := c.hub.expenseService.Create(input, c.user)
	if err != nil {
		fmt.Println("Failed to create expense:", err)
		return
	}

	event, err := c.hub.eventService.FindOne(expense.EventID)
	if err != nil {
		fmt.Println("Failed to get event:", err)
		return
	}

	users, err := c.hub.eventService.GetUsers(expense.EventID)
	if err != nil {
		fmt.Println("Failed to get user users:", err)
	}

	usersTokens := make([]*string, 0)
	for _, user := range users {
		if user.ID != c.user.ID {
			usersTokens = append(usersTokens, user.FirebaseToken)
		}
	}

	notification := messaging.Notification{
		Title: fmt.Sprintf("New expense added to %s", event.Name),
		Body:  fmt.Sprintf("%s, Amount: %.2f", expense.Name, expense.Amount),
	}

	err = firebase.SendNotification(usersTokens, notification)

	c.refreshEventData()
}

func (c *Client) updateExpense(payload json.RawMessage) {
	var input types.UpdateExpenseInput
	if err := json.Unmarshal(payload, &input); err != nil {
		fmt.Println("Failed to unmarshal data:", err)
		return
	}

	_, err := c.hub.expenseService.Update(input.ID, input)
	if err != nil {
		fmt.Println("Failed to update expense:", err)
		return
	}

	c.refreshEventData()
}

func (c *Client) deleteExpense(payload json.RawMessage) {
	var input DeleteExpenseInput
	if err := json.Unmarshal(payload, &input); err != nil {
		fmt.Println("Failed to unmarshal data:", err)
		return
	}

	err := c.hub.expenseService.Delete(input.ID)
	if err != nil {
		fmt.Println("Failed to delete expense:", err)
		return
	}

	c.refreshEventData()
}

func (c *Client) createRefund(payload json.RawMessage) {
	var input types.CreateRefundInput
	if err := json.Unmarshal(payload, &input); err != nil {
		fmt.Println("Failed to unmarshal data:", err)
		return
	}

	eventId, err := strconv.ParseUint(c.room, 10, 32)
	if err != nil {
		fmt.Println("Invalid eventId:", err)
		return
	}

	refund, err := c.hub.refundService.Create(input, c.user, uint(eventId))
	if err != nil {
		fmt.Println("Failed to create expense:", err)
		return
	}

	event, err := c.hub.eventService.FindOne(refund.EventID)
	if err != nil {
		fmt.Println("Failed to get event:", err)
		return
	}

	eventUsers, err := c.hub.eventService.GetUsers(refund.EventID)
	if err != nil {
		fmt.Println("Failed to get event users:", err)
		return
	}

	usersTokens := make([]*string, 0)
	for _, user := range eventUsers {
		if user.ID != c.user.ID {
			usersTokens = append(usersTokens, user.FirebaseToken)
		}
	}

	notification := messaging.Notification{
		Title: fmt.Sprintf("New refund added to %s", event.Name),
		Body:  fmt.Sprintf("%s add a refund of %.2f", c.user.Username, refund.Amount),
	}

	_ = firebase.SendNotification(usersTokens, notification)

	c.refreshEventData()
}

func (c *Client) updateRefund(payload json.RawMessage) {
	var input types.UpdateRefundInput
	if err := json.Unmarshal(payload, &input); err != nil {
		fmt.Println("Failed to unmarshal data:", err)
		return
	}

	_, err := c.hub.refundService.Update(input.ID, input)
	if err != nil {
		fmt.Println("Failed to update refund:", err)
		return
	}

	c.refreshEventData()
}

func (c *Client) deleteRefund(payload json.RawMessage) {
	var input DeleteRefundInput
	if err := json.Unmarshal(payload, &input); err != nil {
		fmt.Println("Failed to unmarshal data:", err)
		return
	}

	err := c.hub.refundService.Delete(input.ID)
	if err != nil {
		fmt.Println("Failed to delete refund:", err)
		return
	}

	c.refreshEventData()
}

func (c *Client) refreshEvent() {
	eventId, err := strconv.ParseUint(c.room, 10, 32)
	if err != nil {
		fmt.Println("Invalid eventId:", err)
		return
	}

	event, err := c.hub.eventService.FindOne(uint(eventId))
	if err != nil {
		fmt.Println("Failed to get event:", err)
		return
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Failed to marshal event:", err)
		return
	}

	message, err := json.Marshal(Message{
		Type:    "event",
		Payload: eventJSON,
	})
	if err != nil {
		fmt.Println("Failed to marshal event message:", err)
		return
	}

	c.hub.rooms[c.room].broadcast <- message

}

func (c *Client) refreshEventData() {
	eventId, err := strconv.ParseUint(c.room, 10, 32)
	if err != nil {
		fmt.Println("Invalid eventId:", err)
		return
	}

	c.refreshExpenses(uint(eventId))
	c.refreshRefunds(uint(eventId))
	c.refreshBalances(uint(eventId))
	c.refreshTransactions(uint(eventId))
	c.refreshUsers(uint(eventId))
}

func (c *Client) refreshExpenses(eventId uint) {
	expenses, err := c.hub.expenseService.FindByEventId(eventId)
	if err != nil {
		fmt.Println("Failed to get expenses:", err)
		return
	}

	expensesJSON, err := json.Marshal(expenses)
	if err != nil {
		fmt.Println("Failed to marshal expenses:", err)
		return
	}

	message, err := json.Marshal(Message{
		Type:    "expenses",
		Payload: expensesJSON,
	})
	if err != nil {
		fmt.Println("Failed to marshal expenses message:", err)
		return
	}

	c.hub.rooms[c.room].broadcast <- message
}

func (c *Client) refreshRefunds(eventId uint) {
	refunds, err := c.hub.refundService.FindByEventId(eventId)
	if err != nil {
		fmt.Println("Failed to get refunds:", err)
		return
	}

	refundsJSON, err := json.Marshal(refunds)
	if err != nil {
		fmt.Println("Failed to marshal refunds:", err)
		return
	}

	message, err := json.Marshal(Message{
		Type:    "refunds",
		Payload: refundsJSON,
	})
	if err != nil {
		fmt.Println("Failed to marshal refunds message:", err)
		return
	}

	c.hub.rooms[c.room].broadcast <- message
}

func (c *Client) refreshBalances(eventId uint) {
	balances, err := c.hub.eventService.GetBalances(eventId)
	if err != nil {
		fmt.Println("Failed to get balances:", err)
		return
	}

	balancesJSON, err := json.Marshal(balances)
	if err != nil {
		fmt.Println("Failed to marshal balances:", err)
		return
	}

	message, err := json.Marshal(Message{
		Type:    "balances",
		Payload: balancesJSON,
	})
	if err != nil {
		fmt.Println("Failed to marshal balances message:", err)
		return
	}

	c.hub.rooms[c.room].broadcast <- message
}

func (c *Client) refreshTransactions(eventId uint) {
	transactions, err := c.hub.eventService.GetTransactions(eventId)
	if err != nil {
		fmt.Println("Failed to create transactions:", err)
		return
	}

	transactionsJSON, err := json.Marshal(transactions)
	if err != nil {
		fmt.Println("Failed to marshal transactions:", err)
		return
	}

	message, err := json.Marshal(Message{
		Type:    "transactions",
		Payload: transactionsJSON,
	})
	if err != nil {
		fmt.Println("Failed to marshal transactions message:", err)
		return
	}

	c.hub.rooms[c.room].broadcast <- message
}

func (c *Client) refreshUsers(eventId uint) {
	users, err := c.hub.eventService.GetUsersWithExpenses(eventId)
	if err != nil {
		fmt.Println("Failed to get users:", err)
		return
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Failed to marshal users:", err)
		return
	}

	message, err := json.Marshal(Message{
		Type:    "users",
		Payload: usersJSON,
	})
	if err != nil {
		fmt.Println("Failed to marshal users message:", err)
		return
	}

	c.hub.rooms[c.room].broadcast <- message
}

package websocket

import "encoding/json"

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type DeleteExpenseInput struct {
	ID uint `json:"id" binding:"required"`
}

type DeleteRefundInput struct {
	ID uint `json:"id" binding:"required"`
}

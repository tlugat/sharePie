package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sharePie-api/internal/auth"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs(hub *Hub, room string, c *gin.Context) {
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

	go client.writePump()
	go client.readPump()
}

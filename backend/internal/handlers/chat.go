package handlers

import (
	"log"
	"net/http"

	"chat-system/backend/internal/hub"
	rdb "chat-system/backend/internal/redis"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ChatHandler(h *hub.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		username := c.GetString("username")
		pin := c.GetString("pin")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("[chat] error upgrade WS: %v", err)
			return
		}

		client := &hub.Client{
			Hub:      h,
			Conn:     conn,
			Send:     make(chan []byte, 256),
			UserID:   userID,
			Username: username,
			Pin:      pin,
		}

		h.Register <- client

		go func() {
			user := rdb.UserFromClient(username, userID, pin)
			if err := rdb.SaveSession(pin, user); err != nil {
				log.Printf("[chat] error guardando sesión: %v", err)
			}
		}()

		go func() {
			pending, err := rdb.PopOffline(pin)
			if err != nil {
				log.Printf("[chat] error leyendo mensajes pendientes: %v", err)
				return
			}
			for _, msg := range pending {
				client.Send <- msg
			}
		}()

		go client.WritePump()
		client.ReadPump()

		rdb.DeleteSession(pin)
	}
}

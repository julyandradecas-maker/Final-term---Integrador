package hub

import (
	"encoding/json"
	"log"
	"time"

	"chat-system/backend/internal/models"
	rdb "chat-system/backend/internal/redis"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   string
	Username string
	Pin      string
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, rawMsg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[client] error leyendo de %s: %v", c.Username, err)
			}
			break
		}

		var in models.WSMessage
		if err := json.Unmarshal(rawMsg, &in); err != nil {
			log.Printf("[client] payload inválido de %s: %v", c.Username, err)
			continue
		}
		if in.To == "" || in.Content == "" {
			continue
		}

		msg := models.Message{
			ID:        uuid.New().String(),
			FromPin:   c.Pin,
			FromUser:  c.Username,
			ToPin:     in.To,
			Content:   in.Content,
			CreatedAt: time.Now(),
		}

		data, err := json.Marshal(msg)
		if err != nil {
			log.Printf("[client] error serializando: %v", err)
			continue
		}

		if !c.Hub.Deliver(in.To, data) {
			if err := rdb.PushOffline(in.To, data); err != nil {
				log.Printf("[client] error guardando mensaje offline: %v", err)
			}
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("[client] error escribiendo a %s: %v", c.Username, err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

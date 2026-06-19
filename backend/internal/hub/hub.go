package hub

import (
	"log"
	"sync"
)

type Hub struct {
	clients    map[string]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.Register:
			h.mu.Lock()
			if h.clients[client.Pin] == nil {
				h.clients[client.Pin] = make(map[*Client]bool)
			}
			h.clients[client.Pin][client] = true
			h.mu.Unlock()
			log.Printf("[hub] %s (pin %s) conectado", client.Username, client.Pin)

		case client := <-h.Unregister:
			h.mu.Lock()
			if conns, ok := h.clients[client.Pin]; ok {
				if _, ok := conns[client]; ok {
					delete(conns, client)
					close(client.Send)
					if len(conns) == 0 {
						delete(h.clients, client.Pin)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("[hub] %s (pin %s) desconectado", client.Username, client.Pin)
		}
	}
}

func (h *Hub) Deliver(toPin string, payload []byte) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	conns := h.clients[toPin]
	if len(conns) == 0 {
		return false
	}
	for client := range conns {
		select {
		case client.Send <- payload:
		default:
			log.Printf("[hub] buffer lleno para pin %s, mensaje descartado", toPin)
		}
	}
	return true
}

func (h *Hub) IsOnline(pin string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.clients[pin]) > 0
}

package websocket

import "log"

var hub *Hub
// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Connection]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Connection

	// Unregister requests from clients.
	unregister chan *Connection

	handler Handler
}

func GetHub() *Hub {
	if hub == nil {
		hub = newHub()
	}
	return hub
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte, 10),
		register:   make(chan *Connection),
		unregister: make(chan *Connection),
		clients:    make(map[*Connection]bool),
	}
}

func (hum *Hub) SetHandler(handler *Handler) {
	hum.handler = *handler
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			log.Println("connection count size : ", len(h.clients))
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

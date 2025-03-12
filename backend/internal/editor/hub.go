package editor

// Hub maintains the set of active clients and broadcasts messages to the
// clients
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// This is the channel to which a client can send messages for broadcasting to others clients
	broadcast chan []byte

	// Register requests from the client
	Register chan *Client

	// Unregister requests from the client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Basically just waits for client messages through channels
// TODO: Return when no more clients are connected.
func (h *Hub) Run() {
	for {
		select {
		// Client registers - just set his value to true
		case client := <-h.Register:
			h.clients[client] = true
			// Client unregisters - set his value to false and
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			// Return if there are not more connected clients
			// if len(h.clients) == 0 {
			// 	return
			// }

			// When a message is broadcast from one of the clients, the hub reads it (from the broadcast channel, obviously)
			// and then sends it over to all the other connected clients who are connected to the same program template
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}
	}
}

package editor

// Hub maintains the set of active clients and broadcasts messages to the
// clients
type Hub struct {
	// Registered clients
	clients map[int64]map[*Client]bool

	// This is the channel to which a client can send messages for broadcasting to others clients
	broadcast chan MessageWrapper

	// Register requests from the client
	Register chan *Client

	// Unregister requests from the client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan MessageWrapper),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[int64]map[*Client]bool),
	}
}

// Basically just waits for client messages through channels
// TODO: Return when no more clients are connected.
func (h *Hub) Run() {
	for {
		select {

		// Client registers - just set his value to true
		case client := <-h.Register:
      if h.clients[client.ProrgamID] == nil {
        h.clients[client.ProrgamID] = make(map[*Client]bool)
      }
			h.clients[client.ProrgamID][client] = true

			// Client unregisters - delete him from the map by ProrgamID and close the send channel
			// to the WritePump
		case client := <-h.unregister:
			if _, ok := h.clients[client.ProrgamID][client]; ok {
				delete(h.clients[client.ProrgamID], client)
				// If the websocket connection is closed, the client unregisters from the ReadPump
				// Then we need to signal the WritePump to quit, so we close the Send channel
				// When it is close, the WritePump sends back a close message to the connection
				// This is required by the WebSockets standard
				close(client.Send)
			}

			// When a message is broadcast from one of the clients, the hub reads it (from the broadcast channel, obviously)
			// and then sends it over to all the other connected clients who are connected to the same program template
		case message := <-h.broadcast:
			for client := range h.clients[message.Id] {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients[message.Id], client)
				}
			}
		}
	}
}

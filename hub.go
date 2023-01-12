package main

import "log"

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Register requests from the clients.
	register chan *Client

	// Rooms.
	rooms map[string]map[*Client]bool

	// Subscription requests from the clients.
	subscribe chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Messages to be broadcasted.
	broadcast chan BroadcastMessage
}

func newHub() *Hub {
	return &Hub{
		clients:   make(map[*Client]bool),
		register:  make(chan *Client),
		rooms:     make(map[string]map[*Client]bool),
		subscribe: make(chan *Client),
		broadcast: make(chan BroadcastMessage),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			room := h.rooms[client.roomID]
			if room != nil {
				if _, ok := room[client]; ok {
					delete(room, client)
					close(client.send)
					if len(room) == 0 {
						// This was last client in the room, delete the room
						delete(h.rooms, client.roomID)
					}
				}
			}
		case client := <-h.subscribe:
			room := h.rooms[client.roomID]
			if room == nil {
				// First client in the room, create a new one
				room = make(map[*Client]bool)
				h.rooms[client.roomID] = room
			}
			room[client] = true
			log.Printf("%+v\n", h.rooms)
		case broadcastMessage := <-h.broadcast:
			room := h.rooms[broadcastMessage.RoomId]
			if room != nil {
				for client := range room {
					select {
					case client.send <- broadcastMessage:
					default:
						close(client.send)
						delete(room, client)
					}
				}
				if len(room) == 0 {
					// The room was emptied while broadcasting to the room.  Delete the room.
					delete(h.rooms, broadcastMessage.RoomId)
				}
			}
		}
	}
}

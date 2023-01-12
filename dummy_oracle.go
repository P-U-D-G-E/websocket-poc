package main

import "time"

type Oracle struct {
	hub *Hub
}

func newOracle(hub *Hub) *Oracle {
	return &Oracle{
		hub: hub,
	}
}

func (o *Oracle) populate() {
	for {
		time.Sleep(3 * time.Second)
		o.hub.broadcast <- BroadcastMessage{RoomId: "1", Data: 1.03}
		o.hub.broadcast <- BroadcastMessage{RoomId: "2", Data: 2.03}
		o.hub.broadcast <- BroadcastMessage{RoomId: "3", Data: 2.45}
	}
}

package main

type Message struct {
	Event         string   `json:"event"`
	SportEventIds []string `json:"sportEventIds"`
}

type BroadcastMessage struct {
	RoomId string  `json:"SportEventId"`
	Data   float64 `json:"data"`
}

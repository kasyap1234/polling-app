package main

import (
	"encoding/json"
	"sync"
	"log"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
	voteCounts map[string]int
}

func NewHub()*Hub{
	return &Hub{
		// list of clients 
		clients: make(map[*Client]bool),
		// used for sending data from hub to client 
		broadcast: make(chan []byte),

		register : make(chan *Client),
		unregister: make(chan *Client),
		voteCounts: make(map[string]int),

	}

}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) HandleVote(message []byte) {
	// decoding the vote 
	var vote struct {
		Option string `json:"Option"` 
	}
	// unmarhsalling it into the vote struct 
	err := json.Unmarshal(message, &vote)
	if err != nil {
		log.Println("Invalid vote format ", err)
		return 
	}
	h.mu.Lock()
	h.voteCounts[vote.Option]++
	updatedCounts := h.voteCounts
	h.mu.Unlock()
	response, err := json.Marshal(updatedCounts)
	if err != nil {
		log.Println("unable to obtain the count value", err)
		return 
	}
	h.broadcast <- response 
}
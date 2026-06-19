package ws

import (
	"encoding/json"
	"log"
)

// Хаб поддерживает набор активных клиентов и рассылает сообщения
// клиентам.

type message struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

type Hub struct {
	//Регистрируем клиента
	clients map[string]*Client

	// Массавая рассылка всем клиентам
	broadcast chan message

	// Регистрировать запросы от клиентов.
	register chan *Client

	// Регистрировать запросы от клиентов.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.userID] = client
		case client := <-h.unregister:
			if c, ok := h.clients[client.userID]; ok && c == client {
				delete(h.clients, client.userID)
				close(client.send)
			}
		case msg := <-h.broadcast:
			if targetClient, ok := h.clients[msg.To]; ok {
				payload, err := json.Marshal(msg)
				if err != nil {
					log.Printf("error marshaling message: %v", err)
					continue
				}
				select {
				case targetClient.send <- payload:
				default:
					close(targetClient.send)
					delete(h.clients, targetClient.userID)
				}
			} else {
				log.Printf("User %s is offline. Message from %s skipped.", msg.To, msg.From)
			}
		}
	}
}

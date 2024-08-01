package config

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	Upgrader  = websocket.Upgrader{}
	RecentHub = NewHub()
)

type Hub struct {
	Clients    map[*websocket.Conn]bool
	Broadcast  chan []byte
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.Register:
			h.Clients[conn] = true
		case conn := <-h.Unregister:
			if _, ok := h.Clients[conn]; ok {
				delete(h.Clients, conn)
				conn.Close()
			}
		case message := <-h.Broadcast:
			for conn := range h.Clients {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					h.Unregister <- conn
					conn.Close()
				}
			}
		}
	}
}

func HandleWebocket(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to ugprade connection", http.StatusInternalServerError)
		return
	}

	RecentHub.Register <- conn

	defer func() {
		RecentHub.Unregister <- conn
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		log.Println("Received message: ", string(message))
	}
}

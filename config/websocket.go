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
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.clients[conn] = true
		case conn := <-h.unregister:
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				conn.Close()
			}
		case message := <-h.broadcast:
			for conn := range h.clients {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					h.unregister <- conn
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

	RecentHub.register <- conn

	defer func() {
		RecentHub.unregister <- conn
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		log.Println("Received message: ", string(message))
	}
}

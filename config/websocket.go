package config

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/repository"
	httputil "github.com/temuka-api-service/pkg/http"
)

var (
	Upgrader  = websocket.Upgrader{}
	RecentHub = NewHub()
)

type Hub struct {
	Clients    map[*websocket.Conn]int
	Broadcast  chan MessagePayload
	Register   chan ClientInfo
	Unregister chan *websocket.Conn
}

type MessagePayload struct {
	ConversationID int    `json:"conversation_id"`
	UserID         int    `json:"user_id"`
	Text           string `json:"text"`
}

type ClientInfo struct {
	Connection     *websocket.Conn
	ConversationID int
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*websocket.Conn]int),
		Broadcast:  make(chan MessagePayload),
		Register:   make(chan ClientInfo),
		Unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.Connection] = client.ConversationID
		case conn := <-h.Unregister:
			if _, ok := h.Clients[conn]; ok {
				delete(h.Clients, conn)
				conn.Close()
			}
		case message := <-h.Broadcast:
			for conn, convID := range h.Clients {
				if convID == message.ConversationID {
					err := conn.WriteJSON(message)
					if err != nil {
						h.Unregister <- conn
						conn.Close()
					}
				}
			}
		}
	}
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	consersationIDstr := r.URL.Query().Get("conversation_id")
	conversationID, err := strconv.Atoi(consersationIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid conversation id"})
		return
	}

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to upgrade connection"})
		return
	}

	RecentHub.Register <- ClientInfo{Connection: conn, ConversationID: conversationID}

	defer func() {
		RecentHub.Unregister <- conn
	}()

	repo := repository.NewMessageRepositoryImpl(GetDBInstance())

	for {
		var message MessagePayload
		err := conn.ReadJSON(&message)
		if err != nil {
			break
		}

		newMessage := model.Message{
			ParticipantID: message.UserID,
			Text:          message.Text,
		}

		err = repo.CreateMessage(Ctx, &newMessage)
		if err != nil {
			log.Println("Failed to save message", err)
		}

		RecentHub.Broadcast <- message
	}

}

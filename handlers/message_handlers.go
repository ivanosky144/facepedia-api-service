package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/temuka-api-service/config"
	"github.com/temuka-api-service/internal"
	"github.com/temuka-api-service/models"
	"gorm.io/gorm"
)

type Client struct {
	Conn        *websocket.Conn
	Message     chan *models.Message
	Participant *models.Participant
	Hub         *internal.Hub
	DB          *gorm.DB
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	var requestBody struct {
		ParticipantID int    `json:"participant_id"`
		Text          string `json:"text"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newMessage := models.Message{
		ParticipantID: requestBody.ParticipantID,
		Text:          requestBody.Text,
	}

	db.Create(&newMessage)

}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		if err := c.Conn.WriteJSON(message); err != nil {
			log.Printf("error: %v", err)
			return
		}
	}
}

func (c *Client) readMessage() {
	defer func() {
		c.Hub.Unregister <- c.Participant
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &models.Message{
			ParticipantID: c.Participant.ID,
			Text:          string(m),
		}

		c.Hub.Broadcast <- msg
	}
}

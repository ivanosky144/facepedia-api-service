package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/temuka-api-service/config"
	"github.com/temuka-api-service/internal"
	"github.com/temuka-api-service/models"
	"gorm.io/gorm"
)

var hub *internal.Hub
var db *gorm.DB

func CreateConversation(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	var requestBody struct {
		Title     string `json:"title"`
		CreatorID int    `json:"creator_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newConversation := models.Conversation{
		Title:  requestBody.Title,
		UserID: requestBody.CreatorID,
	}

	db.Create(&newConversation)

	response := struct {
		Message string              `json:"message"`
		Data    models.Conversation `json:"data"`
	}{
		Message: "Conversation has been created",
		Data:    newConversation,
	}
	respondJSON(w, http.StatusOK, response)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func JoinConversation(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	conversationID, _ := strconv.Atoi(vars["id"])
	userID, _ := strconv.Atoi(r.URL.Query().Get("userId"))

	newParticipant := models.Participant{
		ConversationID: conversationID,
		UserID:         userID,
	}

	if err := db.Create(&newParticipant).Error; err != nil {
		http.Error(w, "Error adding participant", http.StatusInternalServerError)
		return
	}

	client := &Client{
		Conn:        conn,
		Message:     make(chan *models.Message, 10),
		Participant: &newParticipant,
		Hub:         hub,
		DB:          db,
	}

	hub.Register <- &newParticipant

	go client.writeMessage()
	client.readMessage()

}

func GetConversations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	participantIDstr := vars["participant_id"]

	participantID, err := strconv.Atoi(participantIDstr)
	if err != nil {
		http.Error(w, "Invalid participant id", http.StatusBadRequest)
		return
	}

	var conversations []models.Conversation
	err = db.Joins("JOIN participants ON participants.conversation_id = conversations.id").
		Where("participants.id = ?", participantID).
		Preload("Participants").
		Find(&conversations).Error
	if err != nil {
		http.Error(w, "Error retrieving conversations", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string                `json:"message"`
		Data    []models.Conversation `json:"data"`
	}{
		Message: "Timeline posts has been retrieved",
		Data:    conversations,
	}

	respondJSON(w, http.StatusOK, response)
}

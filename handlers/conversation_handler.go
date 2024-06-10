package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/temuka-api-service/models"
)

func CreateConversation(w http.ResponseWriter, r *http.Request) {
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
		Title:     requestBody.Title,
		CreatorID: requestBody.CreatorID,
	}

	db.Create(&newConversation)

	response := struct {
		Message string              `json:"message"`
		Data    models.Conversation `json:"data"`
	}{
		Message: "Conversation has been created",
		Data:    newConversation,
	}
	json.NewEncoder(w).Encode(response)
}

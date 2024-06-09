package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/temuka-api-service/models"
)

func AddComment(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		PostID  int    `json:"post_id"`
		UserID  int    `json:"user_id"`
		Content string `json:"content"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newComment := models.Comment{
		UserID:  requestBody.UserID,
		PostID:  requestBody.PostID,
		Content: requestBody.Content,
	}

	db.Create(&newComment)

	response := struct {
		Message string         `json:"message"`
		Data    models.Comment `json:"data"`
	}{
		Message: "Comment has been added",
		Data:    newComment,
	}

	json.NewEncoder(w).Encode(&response)
}

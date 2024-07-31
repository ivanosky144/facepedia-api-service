package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/temuka-api-service/config"
	"github.com/temuka-api-service/models"
)

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	var requestBody struct {
		UserID int `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var notifications []models.Notification
	if err := db.First(&notifications, "user_id = ?", requestBody.UserID).Error; err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	response := struct {
		Message string                `json:"message"`
		Data    []models.Notification `json:"data"`
	}{
		Message: "Notifications has been retrieved successfully",
		Data:    notifications,
	}

	respondJSON(w, http.StatusOK, response)
}

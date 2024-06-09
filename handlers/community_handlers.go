package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/temuka-api-service/models"
)

func CreateCommunity(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name        string `json:"name"`
		Description string `json:"desc"`
		LogoPicture string `json:"logopicture"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newCommunity := models.Community{
		Name:        requestBody.Name,
		Description: requestBody.Description,
		LogoPicture: requestBody.LogoPicture,
	}

	db.Create(&newCommunity)

	response := struct {
		Message string           `json:"message"`
		Data    models.Community `json:"data"`
	}{
		Message: "Community has been created",
		Data:    newCommunity,
	}
	json.NewEncoder(w).Encode(response)
}

func JoinCommunity(w http.ResponseWriter, r *http.Request) {

}

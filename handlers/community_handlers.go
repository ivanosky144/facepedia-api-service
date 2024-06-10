package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/models"
	"gorm.io/gorm"
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
	vars := mux.Vars(r)
	communityIDstr := vars["id"]

	communityID, err := strconv.Atoi(communityIDstr)
	if err != nil {
		http.Error(w, "Invalid community id", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		UserID int `json:"user_id"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var community models.Community
	if err := db.First(&community, communityID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Community not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving community", http.StatusInternalServerError)
		}
		return
	}

	isMember := false
	for _, memberID := range community.Members {
		if memberID == requestBody.UserID {
			isMember = true
			break
		}
	}

	if isMember {
		http.Error(w, "User already a member of the community", http.StatusBadRequest)
		return
	}

	community.Members = append(community.Members, requestBody.UserID)
	community.MembersCount += 1

	if err := db.Save(&community).Error; err != nil {
		http.Error(w, "Error updating community", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Successfully joined the community",
	}
	json.NewEncoder(w).Encode(response)

}

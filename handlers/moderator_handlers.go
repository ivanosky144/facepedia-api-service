package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/temuka-api-service/config"
	"github.com/temuka-api-service/models"
)

func RequestInvitation(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	var requestBody struct {
		UserID      int    `json:"user_id"`
		Message     string `json:"message"`
		CommunityID int    `json:"community_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.First(&user, "id = ?", requestBody.UserID); err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Invite request has been sent successfully",
	}

	respondJSON(w, http.StatusOK, response)
}

func BanCommunityMember(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	var requestBody struct {
		CommunityMemberID int `json:"community_member_id"`
	}

	var communityMember models.CommunityMember
	if err := db.First(&communityMember, "id = ?", requestBody.CommunityMemberID).Error; err != nil {
		http.Error(w, "Invalid community member id", http.StatusBadRequest)
		return
	}

	communityMember.Banned = false
	if err := db.Save(&communityMember).Error; err != nil {
		http.Error(w, "Error banning user", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Member has been banned",
	}

	respondJSON(w, http.StatusOK, response)
}

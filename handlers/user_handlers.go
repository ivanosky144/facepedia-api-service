package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/config"
	"github.com/temuka-api-service/models"
)

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	var users []models.User

	if err := db.Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(users) == 0 {
		http.Error(w, "No user found", http.StatusNotFound)
		return
	}

	response := struct {
		Message string        `json:"message"`
		Data    []models.User `json:"data"`
	}{
		Message: "User has been created",
		Data:    users,
	}
	respondJSON(w, http.StatusOK, response)
}

func GetUserDetail(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	vars := mux.Vars(r)
	userIDstr := vars["id"]

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	var user models.User
	res := db.First(&user, userID)
	if res.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	db.Find(&user)
	response := struct {
		Message string      `json:"message"`
		Data    models.User `json:"data"`
	}{
		Message: "User detail has been retrieved",
		Data:    user,
	}

	respondJSON(w, http.StatusOK, response)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	var requestBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newUser := models.User{
		Username: requestBody.Username,
		Email:    requestBody.Email,
		Password: requestBody.Password,
	}

	db.Create(&newUser)
	response := struct {
		Message string      `json:"message"`
		Data    models.User `json:"data"`
	}{
		Message: "User has been created",
		Data:    newUser,
	}
	respondJSON(w, http.StatusOK, response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	vars := mux.Vars(r)
	userIDstr := vars["user_id"]

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	var user models.User
	res := db.First(&user, userID)
	if res.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db.Save(&user)
	response := struct {
		Message string      `json:"message"`
		Data    models.User `json:"data"`
	}{
		Message: "User has been updated",
		Data:    user,
	}

	respondJSON(w, http.StatusOK, response)
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	var requestBody struct {
		TargetID      int `json:target_id`
		CurrentUserID int `json:currentuser_id`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var targetUser models.User
	if err := db.First(&targetUser, requestBody.TargetID).Error; err != nil {
		http.Error(w, "Target user not found", http.StatusNotFound)
		return
	}

	if err := db.Exec("INSERT INTO user_follows (user_id, follower_id) VALUES (?, ?)", requestBody.CurrentUserID, requestBody.TargetID).Error; err != nil {
		http.Error(w, "Error following user", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "User has been followed",
	}

	respondJSON(w, http.StatusOK, response)
}

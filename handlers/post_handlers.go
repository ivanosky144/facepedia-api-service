package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/config"
	"github.com/temuka-api-service/models"
	"gorm.io/gorm"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	var requestBody struct {
		Title       string `json:"title"`
		Description string `json:"desc"`
		UserID      int    `json:"user_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newPost := models.Post{
		Title:       requestBody.Title,
		Description: requestBody.Description,
		UserID:      requestBody.UserID,
	}
	db.Create(&newPost)

	response := struct {
		Message string      `json:"message"`
		Data    models.Post `json:"data"`
	}{
		Message: "Post has been created",
		Data:    newPost,
	}
	respondJSON(w, http.StatusOK, response)
}

func GetTimelinePosts(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	vars := mux.Vars(r)
	userIDstr := vars["user_id"]

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	var currentUser models.User
	if err := db.First(&currentUser, "id = ?", userID).Error; err != nil {
		http.Error(w, "Cannot retrieve the data because the user was not found", http.StatusNotFound)
		return
	}

	var userPosts []models.Post
	if err := db.Where("user_id = ?", currentUser.ID).Find(&userPosts).Error; err != nil {
		http.Error(w, "Error retrieving user posts", http.StatusInternalServerError)
		return
	}

	var friendPosts []models.Post
	for _, friendID := range currentUser.Followings {
		var posts []models.Post
		if err := db.Where("user_id = ?", friendID).Find(&posts).Error; err != nil {
			http.Error(w, "Error retrieving friend posts", http.StatusInternalServerError)
			return
		}
		friendPosts = append(friendPosts, posts...)
	}

	timelinePosts := append(userPosts, friendPosts...)

	response := struct {
		Message string        `json:"message"`
		Data    []models.Post `json:"data"`
	}{
		Message: "Timeline posts has been retrieved",
		Data:    timelinePosts,
	}

	respondJSON(w, http.StatusOK, response)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	vars := mux.Vars(r)
	postIDstr := vars["id"]

	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}

	if err := db.Delete(&models.Post{}, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error deleting post", http.StatusInternalServerError)
		}
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Post has been deleted",
	}

	respondJSON(w, http.StatusOK, response)
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	vars := mux.Vars(r)
	postIDstr := vars["id"]

	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
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

	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving post", http.StatusBadRequest)
		}
		return
	}

	alreadyLiked := false
	for _, user := range post.Likes {
		if uint(user.ID) == uint(requestBody.UserID) {
			alreadyLiked = true
			break
		}
	}

	if !alreadyLiked {
		var liker models.User
		if err := db.First(&liker, requestBody.UserID).Error; err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		post.Likes = append(post.Likes, &liker)
		if err := db.Save(&post).Error; err != nil {
			http.Error(w, "Error liking post", http.StatusInternalServerError)
			return
		}

		response := struct {
			Message string `json:"message"`
		}{
			Message: "You have liked this post",
		}

		respondJSON(w, http.StatusOK, response)
	}
}

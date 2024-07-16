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

func AddComment(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	var requestBody struct {
		PostID   int    `json:"post_id"`
		UserID   int    `json:"user_id"`
		ParentID int    `json:parent_id`
		Content  string `json:"content"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newComment := models.Comment{
		UserID:   requestBody.UserID,
		PostID:   requestBody.PostID,
		ParentID: requestBody.ParentID,
		Content:  requestBody.Content,
	}

	db.Create(&newComment)

	response := struct {
		Message string         `json:"message"`
		Data    models.Comment `json:"data"`
	}{
		Message: "Comment has been added",
		Data:    newComment,
	}

	respondJSON(w, http.StatusOK, response)
}

func ShowCommentsByPost(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()
	var requestBody struct {
		PostID int `json:"post_id"`
	}

	var post models.Post
	if err := db.First(&post, "id = ?", requestBody.PostID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Invalid post id", http.StatusBadRequest)
		}
		return
	}

	var postComments []models.Comment
	if err := db.Where("post_id = ?", requestBody.PostID).Find(&postComments).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Comments not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving comments", http.StatusBadRequest)
		}
		return
	}

	response := struct {
		Message string           `json:"message"`
		Data    []models.Comment `json:"data"`
	}{
		Message: "Comments has been retrieved",
		Data:    postComments,
	}

	respondJSON(w, http.StatusOK, response)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	vars := mux.Vars(r)
	commentIDstr := vars["commentId"]

	commentID, err := strconv.Atoi(commentIDstr)
	if err != nil {
		http.Error(w, "Invalid comment id", http.StatusBadRequest)
		return
	}

	if err := db.Delete(&models.Comment{}, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Comment not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error deleting comment", http.StatusBadRequest)
		}
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Comment has been deleted",
	}

	respondJSON(w, http.StatusOK, response)
}

func ShowReplies(w http.ResponseWriter, r *http.Request) {
	db := config.GetDBInstance()

	var requestBody struct {
		ParentID int `json:"parent_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var parent models.Comment
	if err := db.First(&parent, "id = ?", requestBody.ParentID).Error; err != nil {
		http.Error(w, "Invalid parent id", http.StatusBadRequest)
		return
	}

	var replies []models.Comment
	if err := db.Where("parent_id = ?", requestBody.ParentID).Find(&replies).Error; err != nil {
		http.Error(w, "Replies not found", http.StatusNotFound)
		return
	}

	response := struct {
		Message string           `json:"message"`
		Data    []models.Comment `json:"data"`
	}{
		Message: "Replies have been shown",
		Data:    replies,
	}

	respondJSON(w, http.StatusOK, response)
}

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/config"
	"github.com/temuka-api-service/helpers"
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

	var post models.Post
	if err := db.First(&post, "id = ?", requestBody.PostID); err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if post.UserID != requestBody.UserID {
		newCommentNotification := models.Notification{
			UserID:    post.UserID,
			ActorID:   requestBody.UserID,
			CommentID: newComment.ID,
			Type:      "comment",
			Message:   "New comment on your post",
			Read:      false,
		}

		db.Create(&newCommentNotification)
		helpers.PushNotification(newCommentNotification)
	}

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

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
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

	var fetchReplies func(int) ([]models.Comment, error)
	fetchReplies = func(parentID int) ([]models.Comment, error) {
		var comments []models.Comment
		if err := db.Where("parent_id = ?", parentID).Error; err != nil {
			return nil, err
		}

		for i := range comments {
			replies, err := fetchReplies(comments[i].ID)
			if err != nil {
				return nil, err
			}
			comments[i].Replies = replies
		}

		return comments, nil
	}

	replies, err := fetchReplies(requestBody.ParentID)
	if err != nil {
		http.Error(w, "Error retrieving replies", http.StatusBadRequest)
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

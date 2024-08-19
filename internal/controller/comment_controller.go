package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/repository"
	httputil "github.com/temuka-api-service/pkg/http"
)

type CommentController interface {
	AddComment(w http.ResponseWriter, r *http.Request)
	ShowCommentsByPost(w http.ResponseWriter, r *http.Request)
	DeleteComment(w http.ResponseWriter, r *http.Request)
	ShowReplies(w http.ResponseWriter, r *http.Request)
}

type CommentControllerImpl struct {
	CommentRepository      repository.CommentRepository
	PostRepository         repository.PostRepository
	NotificationRepository repository.NotificationRepository
	ReportRepository       repository.ReportRepository
}

func NewCommentController(commentRepo repository.CommentRepository, postRepo repository.PostRepository, notificationRepo repository.NotificationRepository, reportRepo repository.ReportRepository) CommentController {
	return &CommentControllerImpl{
		CommentRepository:      commentRepo,
		PostRepository:         postRepo,
		NotificationRepository: notificationRepo,
		ReportRepository:       reportRepo,
	}
}

func (c *CommentControllerImpl) AddComment(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		PostID   int    `json:"post_id"`
		UserID   int    `json:"user_id"`
		ParentID *int   `json:"parent_id"`
		Content  string `json:"content"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	var parentID *int
	if requestBody.ParentID != nil {
		parentID = requestBody.ParentID
	}

	newComment := model.Comment{
		UserID:   requestBody.UserID,
		PostID:   requestBody.PostID,
		ParentID: parentID,
		Content:  requestBody.Content,
	}

	post, err := c.PostRepository.GetPostDetailByID(context.Background(), newComment.PostID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusNotFound, map[string]string{"error": "Post not found"})
		return
	}

	if err := c.CommentRepository.CreateComment(context.Background(), &newComment); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error creating comment"})
		return
	}

	if post.UserID != requestBody.UserID {
		newCommentNotification := model.Notification{
			UserID:    post.UserID,
			ActorID:   requestBody.UserID,
			PostID:    requestBody.PostID,
			CommentID: newComment.ID,
			Type:      "comment",
			Message:   "New comment on your post",
			Read:      false,
		}

		if err := c.NotificationRepository.CreateNotification(context.Background(), &newCommentNotification); err != nil {
			httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error creating notification"})
			return
		}
	}

	response := struct {
		Message string        `json:"message"`
		Data    model.Comment `json:"data"`
	}{
		Message: "Comment has been added",
		Data:    newComment,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *CommentControllerImpl) ShowCommentsByPost(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		PostID int `json:"post_id"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	comments, err := c.CommentRepository.GetCommentsByPostID(context.Background(), requestBody.PostID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving comments"})
		return
	}

	response := struct {
		Message string          `json:"message"`
		Data    []model.Comment `json:"data"`
	}{
		Message: "Comments have been retrieved",
		Data:    comments,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *CommentControllerImpl) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentIDstr := vars["commentId"]

	commentID, err := strconv.Atoi(commentIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid comment id"})
		return
	}

	if err := c.CommentRepository.DeleteComment(context.Background(), commentID); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error deleting comment"})
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Comment has been deleted",
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *CommentControllerImpl) ShowReplies(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		ParentID int `json:"parent_id"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	var fetchReplies func(parentID int) ([]model.Comment, error)
	fetchReplies = func(parentID int) ([]model.Comment, error) {
		comments, err := c.CommentRepository.GetRepliesByParentID(context.Background(), parentID)
		if err != nil {
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
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving replies"})
		return
	}

	response := struct {
		Message string          `json:"message"`
		Data    []model.Comment `json:"data"`
	}{
		Message: "Replies have been shown",
		Data:    replies,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

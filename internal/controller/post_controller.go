package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/repository"
	httputil "github.com/temuka-api-service/pkg/http"
	"gorm.io/gorm"
)

type PostController interface {
	CreatePost(w http.ResponseWriter, r *http.Request)
	GetPostDetail(w http.ResponseWriter, r *http.Request)
	GetUserPosts(w http.ResponseWriter, r *http.Request)
	UpdatePost(w http.ResponseWriter, r *http.Request)
	DeletePost(w http.ResponseWriter, r *http.Request)
	GetTimelinePosts(w http.ResponseWriter, r *http.Request)
	LikePost(w http.ResponseWriter, r *http.Request)
}

type PostControllerImpl struct {
	PostRepository         repository.PostRepository
	NotificationRepository repository.NotificationRepository
	UserRepository         repository.UserRepository
	ReportRepository       repository.ReportRepository
	CommunityRepository    repository.CommunityRepository
}

func NewPostController(postRepo repository.PostRepository, notificationRepo repository.NotificationRepository, userRepo repository.UserRepository, reportRepo repository.ReportRepository, communityRepo repository.CommunityRepository) PostController {
	return &PostControllerImpl{
		PostRepository:         postRepo,
		NotificationRepository: notificationRepo,
		UserRepository:         userRepo,
		ReportRepository:       reportRepo,
		CommunityRepository:    communityRepo,
	}
}

func (c *PostControllerImpl) CreatePost(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Title       string `json:"title"`
		Description string `json:"desc"`
		UserID      int    `json:"user_id"`
		CommunityID int    `json:"community_id"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	newPost := model.Post{
		Title:       requestBody.Title,
		Description: requestBody.Description,
		UserID:      requestBody.UserID,
	}

	if err := c.PostRepository.CreatePost(context.Background(), &newPost); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error creating post"})
		return
	}

	if err := c.CommunityRepository.UpdateCommunityPostsCount(context.Background(), requestBody.CommunityID); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error updating community posts count"})
		return
	}

	response := struct {
		Message string     `json:"message"`
		Data    model.Post `json:"data"`
	}{
		Message: "Post has been created",
		Data:    newPost,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *PostControllerImpl) GetPostDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIDstr := vars["id"]

	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid post id"})
		return
	}

	post, err := c.PostRepository.GetPostDetailByID(context.Background(), postID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusNotFound, map[string]string{"error": "Post not found"})
		return
	}

	response := struct {
		Message string     `json:"message"`
		Data    model.Post `json:"data"`
	}{
		Message: "Post detail has been retrieved",
		Data:    *post,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *PostControllerImpl) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDstr := vars["user_id"]

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid user id"})
		return
	}

	posts, err := c.PostRepository.GetPostsByUserID(context.Background(), userID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving posts"})
		return
	}

	response := struct {
		Message string       `json:"message"`
		Data    []model.Post `json:"data"`
	}{
		Message: "User posts have been retrieved",
		Data:    posts,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *PostControllerImpl) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIDstr := vars["id"]

	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid post id"})
		return
	}

	var post model.Post
	if err := httputil.ReadRequest(r, &post); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	post.ID = postID
	if err := c.PostRepository.UpdatePost(context.Background(), postID, &post); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error updating post"})
		return
	}

	response := struct {
		Message string     `json:"message"`
		Data    model.Post `json:"data"`
	}{
		Message: "Post has been updated",
		Data:    post,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *PostControllerImpl) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIDstr := vars["id"]

	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid post id"})
		return
	}

	if err := c.PostRepository.DeletePost(context.Background(), postID); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error deleting post"})
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Post has been deleted",
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *PostControllerImpl) GetTimelinePosts(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		UserID int `json:"user_id"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request test"})
		return
	}

	userPosts, err := c.PostRepository.GetPostsByUserID(context.Background(), requestBody.UserID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Error retrieving user posts"})
		return
	}
	userFollowers, err := c.UserRepository.GetFollowers(context.Background(), requestBody.UserID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving followers"})
		return
	}

	var followerPosts []model.Post

	for _, data := range userFollowers {
		if posts, err := c.PostRepository.GetPostsByUserID(context.Background(), data.FollowingID); err == nil {
			followerPosts = append(followerPosts, posts...)
		} else {
			httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving friend posts"})
			return
		}
	}

	var timelinePosts []model.Post
	timelinePosts = append(userPosts, followerPosts...)

	response := struct {
		Message string       `json:"message"`
		Data    []model.Post `json:"data"`
	}{
		Message: "Timeline posts have been retrieved successfully",
		Data:    timelinePosts,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *PostControllerImpl) LikePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIDstr := vars["id"]

	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid post id"})
		return
	}

	var requestBody struct {
		UserID int `json:"user_id"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	post, err := c.PostRepository.GetPostDetailByID(context.Background(), postID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			httputil.WriteResponse(w, http.StatusNotFound, map[string]string{"error": "Post not found"})
		} else {
			httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving post"})
		}
		return
	}

	alreadyLiked := false
	for _, user := range post.Likes {
		if user.ID == requestBody.UserID {
			alreadyLiked = true
			break
		}
	}

	if !alreadyLiked {
		liker, err := c.UserRepository.GetUserByID(context.Background(), requestBody.UserID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				httputil.WriteResponse(w, http.StatusNotFound, map[string]string{"error": "User not found"})
			} else {
				httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving user"})
			}
			return
		}

		post.Likes = append(post.Likes, &model.User{ID: requestBody.UserID})
		if err := c.PostRepository.UpdatePost(context.Background(), postID, post); err != nil {
			httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error liking post"})
			return
		}

		likePostNotification := model.Notification{
			UserID:  post.UserID,
			ActorID: requestBody.UserID,
			PostID:  post.ID,
			Type:    "like",
			Message: liker.Username + " liked your post: " + post.Title,
			Read:    false,
		}
		if err := c.NotificationRepository.CreateNotification(context.Background(), &likePostNotification); err != nil {
			httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error creating notification"})
			return
		}

		response := struct {
			Message string `json:"message"`
		}{
			Message: "You have liked this post",
		}

		httputil.WriteResponse(w, http.StatusOK, response)
	}
}

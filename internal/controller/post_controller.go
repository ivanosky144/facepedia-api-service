package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"time"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/repository"
	httputil "github.com/temuka-api-service/pkg/http"
	"github.com/temuka-api-service/pkg/redis"
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
	CommentRepository      repository.CommentRepository
}

func NewPostController(postRepo repository.PostRepository, notificationRepo repository.NotificationRepository, userRepo repository.UserRepository, reportRepo repository.ReportRepository, communityRepo repository.CommunityRepository, commentRepo repository.CommentRepository) PostController {
	return &PostControllerImpl{
		PostRepository:         postRepo,
		NotificationRepository: notificationRepo,
		UserRepository:         userRepo,
		ReportRepository:       reportRepo,
		CommunityRepository:    communityRepo,
		CommentRepository:      commentRepo,
	}
}

func (c *PostControllerImpl) CreatePost(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Title       string `json:"title"`
		Description string `json:"description"`
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

	user, err := c.UserRepository.GetUserByID(context.Background(), post.UserID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	comments, err := c.CommentRepository.GetCommentsByPostID(context.Background(), postID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving comments"})
		return
	}

	type UserData struct {
		Username       string `json:"Username"`
		ProfilePicture string `json:"ProfilePicture"`
	}

	type Comment struct {
		ID        int       `json:"ID"`
		Content   string    `json:"Content"`
		Username  string    `json:"Username"`
		UserPhoto string    `json:"Userphoto"`
		Votes     int       `json:"Votes"`
		CreatedAt time.Time `json:"CreatedAt"`
		UpdatedAt time.Time `json:"UpdatedAt"`
	}

	postComments := make([]Comment, 0, len(comments))
	for _, comment := range comments {
		commentUser, err := c.UserRepository.GetUserByID(context.Background(), comment.UserID)
		if err != nil {
			httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving comment user"})
			return
		}

		postComments = append(postComments, Comment{
			ID:        comment.ID,
			Content:   comment.Content,
			Username:  commentUser.Username,
			UserPhoto: commentUser.ProfilePicture,
			Votes:     len(comment.Votes),
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		})
	}

	userData := UserData{
		Username:       user.Username,
		ProfilePicture: user.ProfilePicture,
	}

	type ResponseData struct {
		User     UserData   `json:"user"`
		Post     model.Post `json:"post"`
		Comments []Comment  `json:"comments"`
	}

	response := struct {
		Message string       `json:"message"`
		Data    ResponseData `json:"data"`
	}{
		Message: "Post detail has been retrieved",
		Data: ResponseData{
			User:     userData,
			Post:     *post,
			Comments: postComments,
		},
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

	var requestBody struct {
		UserID      int    `json:"user_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	updatedPost := model.Post{
		UserID:      requestBody.UserID,
		Title:       requestBody.Title,
		Description: requestBody.Description,
	}

	post.ID = postID
	if err := c.PostRepository.UpdatePost(context.Background(), postID, &updatedPost); err != nil {
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
	vars := mux.Vars(r)
	userIDstr := vars["user_id"]

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid user id"})
		return
	}

	cacheKey := fmt.Sprintf("timeline_posts_user_%d", userID)

	var cachedResponse struct {
		Message string       `json:"message"`
		Data    []model.Post `json:"data"`
	}

	if err := redis.GetCache(cacheKey, &cachedResponse); err == nil {
		log.Printf("Cache hit for user %d", userID)
		httputil.WriteResponse(w, http.StatusOK, cachedResponse)
		return
	} else if err != nil {
		log.Printf("Cache miss for user %d", userID)
	}

	userPosts, err := c.PostRepository.GetPostsByUserID(context.Background(), userID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Error retrieving user posts"})
		return
	}
	userFollowers, err := c.UserRepository.GetFollowers(context.Background(), userID)
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

	if err := redis.SetCache(cacheKey, response, 10*time.Minute); err != nil {
		log.Print("Error caching timeline posts:")
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

package router

import (
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/handlers"
)

func Routes() *mux.Router {
	router := mux.NewRouter()

	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/login", handlers.Login).Methods("POST")
	authRouter.HandleFunc("/register", handlers.Register).Methods("POST")
	authRouter.HandleFunc("/resetPassword/{id}", handlers.ResetPassword).Methods("POST")

	userRouter := router.PathPrefix("/api/user").Subrouter()
	userRouter.HandleFunc("/", handlers.CreateUser).Methods("POST")
	userRouter.HandleFunc("/{id}", handlers.UpdateUser).Methods("PUT")
	userRouter.HandleFunc("/search", handlers.SearchUsers).Methods("GET")
	userRouter.HandleFunc("/follow", handlers.FollowUser).Methods("POST")
	userRouter.HandleFunc("/followers", handlers.GetFollowers).Methods("GET")
	userRouter.HandleFunc("/{id}", handlers.GetUserDetail).Methods("GET")

	postRouter := router.PathPrefix("/api/post").Subrouter()
	postRouter.HandleFunc("/create", handlers.CreatePost).Methods("POST")
	postRouter.HandleFunc("/timeline/{userId}", handlers.GetTimelinePosts).Methods("GET")
	postRouter.HandleFunc("/like/{id}", handlers.LikePost).Methods("PUT")
	postRouter.HandleFunc("/{id}", handlers.DeletePost).Methods("DELETE")
	postRouter.HandleFunc("/{id}", handlers.GetPostDetail).Methods("GET")
	postRouter.HandleFunc("/{id}", handlers.UpdatePost).Methods("PUT")

	commentRouter := router.PathPrefix("/api/comment").Subrouter()
	commentRouter.HandleFunc("/add", handlers.AddComment).Methods("POST")
	commentRouter.HandleFunc("/replies", handlers.ShowReplies).Methods("GET")
	commentRouter.HandleFunc("/{commentId}", handlers.DeleteComment).Methods("DELETE")
	commentRouter.HandleFunc("/show", handlers.ShowCommentsByPost).Methods("GET")

	communityRouter := router.PathPrefix("/api/community").Subrouter()
	communityRouter.HandleFunc("/", handlers.CreateCommunity).Methods("POST")
	communityRouter.HandleFunc("/join/{community_id}", handlers.JoinCommunity).Methods("POST")

	moderatorRouter := router.PathPrefix("/api/moderator").Subrouter()
	moderatorRouter.HandleFunc("/invite", handlers.RequestInvitation).Methods("POST")
	moderatorRouter.HandleFunc("/ban", handlers.BanCommunityMember).Methods("POST")

	fileRouter := router.PathPrefix("/api/file").Subrouter()
	fileRouter.HandleFunc("/file", handlers.UploadFile).Methods("POST")

	notificationRouter := router.PathPrefix("/api/notification").Subrouter()
	notificationRouter.HandleFunc("/list", handlers.GetNotifications).Methods("POST")

	return router
}

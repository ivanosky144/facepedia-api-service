package router

import (
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/internal/controller"
	"github.com/temuka-api-service/internal/repository"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()

	// Init repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	communityRepo := repository.NewCommunityRepository(db)
	moderatorRepo := repository.NewModeratorRepository(db)

	// Init controllers
	authController := controller.NewAuthController(userRepo)
	userController := controller.NewUserController(userRepo)
	postController := controller.NewPostController(postRepo, notificationRepo, userRepo)
	communityController := controller.NewCommunityController(communityRepo)
	commentController := controller.NewCommentController(commentRepo, postRepo, notificationRepo)
	notificationController := controller.NewNotificationController(notificationRepo)
	moderatorController := controller.NewModeratorController(moderatorRepo, notificationRepo)
	fileUploadController := controller.NewFileUploadController("uploads")

	// Init routers
	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/login", authController.Login).Methods("POST")
	authRouter.HandleFunc("/register", authController.Register).Methods("POST")
	authRouter.HandleFunc("/resetPassword/{id}", authController.ResetPassword).Methods("POST")

	userRouter := router.PathPrefix("/api/user").Subrouter()
	userRouter.HandleFunc("/", userController.CreateUser).Methods("POST")
	userRouter.HandleFunc("/{id}", userController.UpdateUser).Methods("PUT")
	userRouter.HandleFunc("/search", userController.SearchUsers).Methods("GET")
	userRouter.HandleFunc("/follow", userController.FollowUser).Methods("POST")
	userRouter.HandleFunc("/followers", userController.GetFollowers).Methods("GET")
	userRouter.HandleFunc("/{id}", userController.GetUserDetail).Methods("GET")

	postRouter := router.PathPrefix("/api/post").Subrouter()
	postRouter.HandleFunc("/create", postController.CreatePost).Methods("POST")
	postRouter.HandleFunc("/timeline", postController.GetTimelinePosts).Methods("GET")
	postRouter.HandleFunc("/self/{user_id}", postController.GetUserPosts).Methods("GET")
	postRouter.HandleFunc("/like/{id}", postController.LikePost).Methods("PUT")
	postRouter.HandleFunc("/{id}", postController.DeletePost).Methods("DELETE")
	postRouter.HandleFunc("/{id}", postController.GetPostDetail).Methods("GET")
	postRouter.HandleFunc("/{id}", postController.UpdatePost).Methods("PUT")

	commentRouter := router.PathPrefix("/api/comment").Subrouter()
	commentRouter.HandleFunc("/add", commentController.AddComment).Methods("POST")
	commentRouter.HandleFunc("/replies", commentController.ShowReplies).Methods("GET")
	commentRouter.HandleFunc("/{commentId}", commentController.DeleteComment).Methods("DELETE")
	commentRouter.HandleFunc("/show", commentController.ShowCommentsByPost).Methods("GET")

	communityRouter := router.PathPrefix("/api/community").Subrouter()
	communityRouter.HandleFunc("/", communityController.CreateCommunity).Methods("POST")
	communityRouter.HandleFunc("/join/{community_id}", communityController.JoinCommunity).Methods("POST")

	fileRouter := router.PathPrefix("/api/file").Subrouter()
	fileRouter.HandleFunc("/", fileUploadController.Upload).Methods("POST")

	notificationRouter := router.PathPrefix("/api/notification").Subrouter()
	notificationRouter.HandleFunc("/list", notificationController.GetNotificationsByUser).Methods("POST")

	moderatorRouter := router.PathPrefix("/api/moderator").Subrouter()
	moderatorRouter.HandleFunc("/send", moderatorController.SendModeratorRequest).Methods("POST")
	moderatorRouter.HandleFunc("/id", moderatorController.RemoveModerator).Methods("DELETE")

	return router
}

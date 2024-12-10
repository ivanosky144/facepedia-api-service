package router

import (
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/internal/controller"
	"github.com/temuka-api-service/internal/repository"
	"github.com/temuka-api-service/middleware"
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
	reportRepo := repository.NewReportRepository(db)
	universityRepo := repository.NewUniversityRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	locationRepo := repository.NewLocationRepository(db)
	conversationRepo := repository.NewConversationRepository(db)

	// Init controllers
	authController := controller.NewAuthController(userRepo)
	userController := controller.NewUserController(userRepo)
	postController := controller.NewPostController(postRepo, notificationRepo, userRepo, reportRepo, communityRepo)
	communityController := controller.NewCommunityController(communityRepo)
	commentController := controller.NewCommentController(commentRepo, postRepo, notificationRepo, reportRepo)
	notificationController := controller.NewNotificationController(notificationRepo)
	moderatorController := controller.NewModeratorController(moderatorRepo, notificationRepo)
	reportController := controller.NewReportController(reportRepo)
	universityController := controller.NewUniversityController(universityRepo, reviewRepo)
	locationController := controller.NewLocationController(locationRepo)
	conversationController := controller.NewConversationController(conversationRepo, userRepo)
	fileUploadController := controller.NewFileUploadController("uploads")

	// Init routers
	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/login", authController.Login).Methods("POST")
	authRouter.HandleFunc("/register", authController.Register).Methods("POST")
	authRouter.HandleFunc("/resetPassword/{id}", authController.ResetPassword).Methods("POST")

	userRouter := router.PathPrefix("/api/user").Subrouter()
	userRouter.Use(middleware.CheckAuth)
	userRouter.HandleFunc("", userController.CreateUser).Methods("POST")
	userRouter.HandleFunc("/{id}", userController.UpdateUser).Methods("PUT")
	userRouter.HandleFunc("/search", userController.SearchUsers).Methods("GET")
	userRouter.HandleFunc("/follow", userController.FollowUser).Methods("POST")
	userRouter.HandleFunc("/followers", userController.GetFollowers).Methods("GET")
	userRouter.HandleFunc("/{id}", userController.GetUserDetail).Methods("GET")

	postRouter := router.PathPrefix("/api/post").Subrouter()
	postRouter.Use(middleware.CheckAuth)
	postRouter.HandleFunc("", postController.CreatePost).Methods("POST")
	postRouter.HandleFunc("/timeline/{user_id}", postController.GetTimelinePosts).Methods("GET")
	postRouter.HandleFunc("/{user_id}", postController.GetUserPosts).Methods("GET")
	postRouter.HandleFunc("/like/{id}", postController.LikePost).Methods("PUT")
	postRouter.HandleFunc("/{id}", postController.DeletePost).Methods("DELETE")
	postRouter.HandleFunc("/{id}", postController.GetPostDetail).Methods("GET")
	postRouter.HandleFunc("/{id}", postController.UpdatePost).Methods("PUT")

	commentRouter := router.PathPrefix("/api/comment").Subrouter()
	commentRouter.Use(middleware.CheckAuth)
	commentRouter.HandleFunc("", commentController.AddComment).Methods("POST")
	commentRouter.HandleFunc("/replies", commentController.ShowReplies).Methods("GET")
	commentRouter.HandleFunc("/{commentId}", commentController.DeleteComment).Methods("DELETE")
	commentRouter.HandleFunc("/show", commentController.ShowCommentsByPost).Methods("GET")

	communityRouter := router.PathPrefix("/api/community").Subrouter()
	communityRouter.Use(middleware.CheckAuth)
	communityRouter.HandleFunc("", communityController.CreateCommunity).Methods("POST")
	communityRouter.HandleFunc("", communityController.GetCommunities).Methods("GET")
	communityRouter.HandleFunc("/join/{community_id}", communityController.JoinCommunity).Methods("POST")
	communityRouter.HandleFunc("/post/{id}", communityController.GetCommunityPosts).Methods("GET")
	communityRouter.HandleFunc("/user", communityController.GetUserJoinedCommunities).Methods("POST")
	communityRouter.HandleFunc("/{id}", communityController.GetCommunityDetail).Methods("GET")
	communityRouter.HandleFunc("/{id}", communityController.DeleteCommunity).Methods("DELETE")
	communityRouter.HandleFunc("/{id}", communityController.UpdateCommunity).Methods("PUT")

	fileRouter := router.PathPrefix("/api/file").Subrouter()
	fileRouter.Use(middleware.CheckAuth)
	fileRouter.HandleFunc("", fileUploadController.Upload).Methods("POST")

	notificationRouter := router.PathPrefix("/api/notification").Subrouter()
	notificationRouter.HandleFunc("/list/{user_id}", notificationController.GetNotificationsByUser).Methods("GET")

	moderatorRouter := router.PathPrefix("/api/moderator").Subrouter()
	moderatorRouter.Use(middleware.CheckAuth)
	moderatorRouter.HandleFunc("/send", moderatorController.SendModeratorRequest).Methods("POST")
	moderatorRouter.HandleFunc("/{id}", moderatorController.RemoveModerator).Methods("DELETE")

	reportRouter := router.PathPrefix("/api/report").Subrouter()
	reportRouter.Use(middleware.CheckAuth)
	reportRouter.HandleFunc("", reportController.CreateReport).Methods("POST")
	reportRouter.HandleFunc("/{id}", reportController.DeleteReport).Methods("DELETE")

	universityRouter := router.PathPrefix("/api/university").Subrouter()
	universityRouter.Use(middleware.CheckAuth)
	universityRouter.HandleFunc("", universityController.AddUniversity).Methods("POST")
	universityRouter.HandleFunc("/{id}", universityController.UpdateUniversity).Methods("PUT")
	universityRouter.HandleFunc("/{slug}", universityController.GetUniversityDetail).Methods("GET")
	universityRouter.HandleFunc("", universityController.GetUniversities).Methods("GET")
	universityRouter.HandleFunc("/review", universityController.AddReview).Methods("POST")
	universityRouter.HandleFunc("/review/university_id", universityController.GetUniversityReviews).Methods("GET")

	locationRouter := router.PathPrefix("/api/location").Subrouter()
	locationRouter.Use(middleware.CheckAuth)
	locationRouter.HandleFunc("", locationController.AddLocation).Methods("POST")
	locationRouter.HandleFunc("", locationController.GetLocations).Methods("GET")
	locationRouter.HandleFunc("/{id}", locationController.UpdateLocation).Methods("PUT")

	conversationRouter := router.PathPrefix("/api/conversation").Subrouter()
	conversationRouter.Use(middleware.CheckAuth)
	conversationRouter.HandleFunc("", conversationController.AddConversation).Methods("POST")
	conversationRouter.HandleFunc("", conversationController.DeleteConversation).Methods("DELETE")

	return router
}

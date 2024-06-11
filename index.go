package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/config"
	"github.com/temuka-api-service/models"
	"github.com/temuka-api-service/routes"
)

func EnableCors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func main() {
	config.OpenConnection()
	if err := config.Database.AutoMigrate(&models.User{}, &models.Community{}, &models.Post{}, &models.Conversation{}, &models.Comment{}, &models.CommunityMember{}, &models.CommunityPost{}, &models.Moderator{}, &models.Participant{}); err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}
	log.Println("Auto-migration completed.")
	router := mux.NewRouter()

	router.PathPrefix("/auth").Handler(routes.AuthRoutes())
	router.PathPrefix("/user").Handler(routes.UserRoutes())
	router.PathPrefix("/post").Handler(routes.PostRoutes())
	router.PathPrefix("/conversation").Handler(routes.ConversationRoutes())
	router.PathPrefix("/comment").Handler(routes.CommentRoutes())
	router.PathPrefix("/community").Handler(routes.CommunityRoutes())

	http.Handle("/", router)

	log.Println("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

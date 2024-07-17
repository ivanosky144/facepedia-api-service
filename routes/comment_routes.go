package routes

import (
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/handlers"
)

func CommentRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/add", handlers.AddComment).Methods("POST")
	r.HandleFunc("/replies", handlers.ShowReplies).Methods("GET")
	r.HandleFunc("/{commentId}", handlers.DeleteComment).Methods("DELETE")
	r.HandleFunc("/show", handlers.ShowCommentsByPost).Methods("GET")

	return r
}

package routes

import (
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/handlers"
)

func PostRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/create", handlers.CreatePost).Methods("POST")
	r.HandleFunc("/timeline/{userId}", handlers.GetTimelinePosts).Methods("GET")
	r.HandleFunc("/like/{id}", handlers.LikePost).Methods("PUT")
	r.HandleFunc("/{id}", handlers.DeletePost).Methods("DELETE")
	r.HandleFunc("/{id}", handlers.GetPostDetail).Methods("GET")
	r.HandleFunc("/{id}", handlers.UpdatePost).Methods("PUT")

	return r
}

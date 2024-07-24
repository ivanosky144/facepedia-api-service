package routes

import (
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/handlers"
)

func UserRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/search", handlers.SearchUsers).Methods("GET")
	r.HandleFunc("/follow", handlers.FollowUser).Methods("POST")
	r.HandleFunc("/followers", handlers.GetFollowers).Methods("GET")
	r.HandleFunc("/{id}", handlers.GetUserDetail).Methods("GET")
	r.HandleFunc("/{id}", handlers.UpdateUser).Methods("PUT")

	return r
}

package routes

import (
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/handlers"
)

func UserRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/", handlers.GetUserDetail).Methods("GET")
	r.HandleFunc("/{id}", handlers.SearchUsers).Methods("GET")

	return r
}

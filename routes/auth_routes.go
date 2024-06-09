package routes

import (
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/handlers"
)

func AuthRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/resetPassword/{id}", handlers.ResetPassword).Methods("POST")

	return r
}

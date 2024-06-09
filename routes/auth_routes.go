package routes

import (
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/controller"
)

func AuthRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/login", controller.Login).Methods("POST")
	r.HandleFunc("/register", controller.Register).Methods("POST")
	r.HandleFunc("/resetPassword/{id}", controller.ResetPassword).Methods("POST")

	return r
}

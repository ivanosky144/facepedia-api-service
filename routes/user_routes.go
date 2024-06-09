package routes

import (
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/controller"
)

func UserRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", controller.CreateUser).Methods("POST")
	r.HandleFunc("/", controller.GetUserDetail).Methods("GET")
	r.HandleFunc("/{id}", controller.SearchUsers).Methods("GET")

	return r
}

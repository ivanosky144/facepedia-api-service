package routes

import (
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/handlers"
)

func NotificationRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/list", handlers.GetNotifications).Methods("POST")

	return r
}

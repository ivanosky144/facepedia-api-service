package routes

import (
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/handlers"
)

func ConversationRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.CreateConversation).Methods("POST")

	return r
}

package routes

import (
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/handlers"
)

func ModeratorRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/invite", handlers.RequestInvitation).Methods("POST")
	r.HandleFunc("/ban", handlers.BanCommunityMember).Methods("POST")

	return r
}

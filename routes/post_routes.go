package routes

import (
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/controller"
)

func PostRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", controller.CreatePost).Methods("POST")
	r.HandleFunc("/timeline", controller.GetTimelinePosts).Methods("GET")
	r.HandleFunc("/like/{id}", controller.LikePost).Methods("PUT")
	r.HandleFunc("/{id}", controller.DeletePost).Methods("DELETE")

	return r
}

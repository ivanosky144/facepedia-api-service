package routes

import (
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/handlers"
)

func FileRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/file", handlers.UploadFile).Methods("POST")

	return r
}

package main

import (
	"log"
	"net/http"

	router "github.com/temuka-api-service/api"
	"github.com/temuka-api-service/config"
	"github.com/temuka-api-service/internal/model"
	"gorm.io/gorm"
)

func EnableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	config.OpenConnection()
	var db *gorm.DB = config.GetDBInstance()

	if config.Database == nil {
		log.Fatal("Database connection is nil")
	}
	if err := config.Database.AutoMigrate(&model.User{}, &model.Community{}, &model.Post{}, &model.Conversation{}, &model.Comment{}, &model.CommunityMember{}, &model.CommunityPost{}, &model.Moderator{}, &model.Participant{}, &model.UserFollow{}, &model.Notification{}, &model.Report{}, &model.Review{}, &model.Location{}, &model.University{}); err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	log.Printf("Database : %v", db)
	log.Println("Auto-migration completed.")

	router := router.Routes(db)
	router.Use(EnableCors)

	http.Handle("/", router)

	log.Println("Server is listening on port 3200")
	log.Fatal(http.ListenAndServe("localhost:3200", nil))
}

package main

import (
	"log"
	"net/http"

	router "github.com/temuka-api-service/api"
	"github.com/temuka-api-service/config"
	"github.com/temuka-api-service/models"
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
	if err := config.Database.AutoMigrate(&models.User{}, &models.Community{}, &models.Post{}, &models.Conversation{}, &models.Comment{}, &models.CommunityMember{}, &models.CommunityPost{}, &models.Moderator{}, &models.Participant{}, &models.UserFollow{}, &models.Notification{}); err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}
	log.Printf("Database : %v", db)
	log.Println("Auto-migration completed.")

	router := router.Routes(db)
	router.Use(EnableCors)

	http.Handle("/", router)

	log.Println("Server is listening on port 3200")
	log.Fatal(http.ListenAndServe("localhost:3200", nil))
}

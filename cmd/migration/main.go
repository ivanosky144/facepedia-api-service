package main

import (
	"log"

	"github.com/temuka-api-service/config"
	"github.com/temuka-api-service/internal/model"
)

func main() {
	config.OpenConnection()

	if config.Database == nil {
		log.Fatal("Database connection is nil")
	}

	if err := config.Database.AutoMigrate(
		&model.User{},
		&model.Community{},
		&model.Post{},
		&model.Conversation{},
		&model.Comment{},
		&model.CommunityMember{},
		&model.CommunityPost{},
		&model.Moderator{},
		&model.Participant{},
		&model.UserFollow{},
		&model.Notification{},
		&model.Report{},
		&model.Location{},
		&model.University{},
		&model.Review{},
		&model.Major{},
		&model.MajorReview{},
	); err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	log.Println("Database migration completed successfully.")
}

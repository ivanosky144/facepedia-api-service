package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Database *gorm.DB
	err      error
)

func OpenConnection() {

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading env file: %v", err)
	}

	pgHost := os.Getenv("PG_HOST")
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")
	pgPort := os.Getenv("PG_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		pgHost, pgUser, pgPassword, dbName, pgPort,
	)

	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connection established successfully")

}

func GetDBInstance() *gorm.DB {
	return Database
}

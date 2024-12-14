package config

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() {

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading env file: %v", err)
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisUser := os.Getenv("REDIS_USER")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Username: redisUser,
		Password: redisPassword,
		DB:       0,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis")
	}

	log.Println("Connected to Redis")
}

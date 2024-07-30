package config

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis")
	}

	log.Println("Connected to Redis")
}

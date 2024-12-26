package queue

import (
	"context"
	"log"

	"github.com/temuka-api-service/pkg/redis"
)

func StartListening(ctx context.Context) {
	go subscribeToQueue(ctx, "content-interaction-queue")
	go subscribeToQueue(ctx, "education-interaction-queue")
}

func subscribeToQueue(ctx context.Context, queueName string) {
	log.Printf("Subscribing to queue: %s", queueName)

	handler := func(msg *redis.Message) {
		log.Printf("Received message from %s: UserID=%s, Action")
	}

	redis.Subscribe(ctx, queueName, handler)
}

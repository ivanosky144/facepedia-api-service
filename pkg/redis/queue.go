package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/temuka-api-service/config"
)

type Message struct {
	UserID string `json:"user_id"`
	Action string `json:"action"`
}

func PublishMessage(ctx context.Context, queueName string, message Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = config.RedisClient.Publish(ctx, queueName, data).Err()
	if err != nil {
		log.Printf("Error publishing to queue %s: %v", queueName, err)
	}
	return err
}

func Subscribe(ctx context.Context, queueName string, handler func(*Message)) {
	sub := config.RedisClient.Subscribe(ctx, queueName)
	defer sub.Close()

	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("Error receiving message from queue %s: %v", queueName, err)
			continue
		}

		var message Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Printf("Error decoding message from queue %s: %v", queueName, err)
			continue
		}

		handler(&message)
	}
}

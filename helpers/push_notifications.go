package helpers

import (
	"encoding/json"
	"log"

	"github.com/temuka-api-service/config"
	"github.com/temuka-api-service/models"
)

func PushNotification(notification models.Notification) {
	msg, err := json.Marshal(notification)
	if err != nil {
		log.Println("Error marshaling notification")
		return
	}

	config.RedisClient.Publish(config.Ctx, "notifications", msg)
}

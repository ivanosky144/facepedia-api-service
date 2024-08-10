package redis

import (
	"encoding/json"
	"time"

	"github.com/temuka-api-service/config"
)

func SetCache(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return config.RedisClient.Set(config.Ctx, key, data, expiration).Err()
}

func GetCache(key string, dest interface{}) error {
	data, err := config.RedisClient.Get(config.Ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

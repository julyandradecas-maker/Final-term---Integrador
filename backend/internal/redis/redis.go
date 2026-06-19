package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"chat-system/backend/internal/config"
	"chat-system/backend/internal/models"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init(cfg *config.Config) error {
	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return fmt.Errorf("URL de Redis inválida: %w", err)
	}

	Client = redis.NewClient(opts)

	ctx := context.Background()
	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("no se pudo conectar a Redis en %s: %w", cfg.RedisURL, err)
	}

	log.Printf("[redis] conectado en %s", cfg.RedisURL)
	return nil
}

func SaveSession(pin string, user models.User) error {
	ctx := context.Background()
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error serializando usuario: %w", err)
	}
	key := fmt.Sprintf("session:%s", pin)
	return Client.Set(ctx, key, data, 24*time.Hour).Err()
}

func DeleteSession(pin string) error {
	ctx := context.Background()
	return Client.Del(ctx, fmt.Sprintf("session:%s", pin)).Err()
}

func PushOffline(pin string, payload []byte) error {
	ctx := context.Background()
	key := fmt.Sprintf("offline:%s", pin)
	if err := Client.RPush(ctx, key, payload).Err(); err != nil {
		return err
	}
	return Client.Expire(ctx, key, 7*24*time.Hour).Err()
}

func PopOffline(pin string) ([][]byte, error) {
	ctx := context.Background()
	key := fmt.Sprintf("offline:%s", pin)

	items, err := Client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	Client.Del(ctx, key)

	out := make([][]byte, 0, len(items))
	for _, it := range items {
		out = append(out, []byte(it))
	}
	return out, nil
}

func UserFromClient(username, userID, pin string) models.User {
	return models.User{ID: userID, Username: username, Pin: pin}
}

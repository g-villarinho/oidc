package redis

import (
	"context"
	"fmt"

	"github.com/g-villarinho/oidc-server/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(config *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connect to redis: %w", err)
	}

	return client, nil
}

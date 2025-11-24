package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *cache {
	return &cache{
		client: client,
	}
}

func (c *cache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("key not found: %s", key)
		}
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return val, nil
}

func (c *cache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

func (c *cache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

func (c *cache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}
	return count > 0, nil
}

func (c *cache) SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	success, err := c.client.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		return false, fmt.Errorf("failed to set nx key %s: %w", key, err)
	}
	return success, nil
}

func (c *cache) GetDel(ctx context.Context, key string) (string, error) {
	val, err := c.client.GetDel(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("key not found: %s", key)
		}
		return "", fmt.Errorf("failed to get del key %s: %w", key, err)
	}
	return val, nil
}

func (c *cache) MGet(ctx context.Context, keys ...string) ([]string, error) {
	vals, err := c.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to mget keys: %w", err)
	}

	result := make([]string, len(vals))
	for i, val := range vals {
		if val != nil {
			result[i] = val.(string)
		}
	}

	return result, nil
}

func (c *cache) MSet(ctx context.Context, pairs map[string]string, ttl time.Duration) error {
	pipe := c.client.Pipeline()

	for key, value := range pairs {
		pipe.Set(ctx, key, value, ttl)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to mset: %w", err)
	}

	return nil
}

func (c *cache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	err := c.client.Expire(ctx, key, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to expire key %s: %w", key, err)
	}
	return nil
}

func (c *cache) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get ttl for key %s: %w", key, err)
	}

	if ttl == -2*time.Second {
		return 0, fmt.Errorf("key not found: %s", key)
	}

	if ttl == -1*time.Second {
		return 0, nil
	}

	return ttl, nil
}

func (c *cache) Increment(ctx context.Context, key string) (int64, error) {
	val, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key %s: %w", key, err)
	}
	return val, nil
}

func (c *cache) Decrement(ctx context.Context, key string) (int64, error) {
	val, err := c.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement key %s: %w", key, err)
	}
	return val, nil
}

func (c *cache) FlushAll(ctx context.Context) error {
	err := c.client.FlushAll(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to flush all: %w", err)
	}
	return nil
}

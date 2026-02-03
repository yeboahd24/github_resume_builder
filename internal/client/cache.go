package client

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheClient struct {
	client  *redis.Client
	enabled bool
}

func NewCacheClient(addr, password string, db int, enabled bool) *CacheClient {
	if !enabled {
		return &CacheClient{enabled: false}
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &CacheClient{
		client:  client,
		enabled: true,
	}
}

func (c *CacheClient) Get(ctx context.Context, key string) (string, error) {
	if !c.enabled {
		return "", redis.Nil
	}
	return c.client.Get(ctx, key).Result()
}

func (c *CacheClient) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if !c.enabled {
		return nil
	}
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *CacheClient) Delete(ctx context.Context, key string) error {
	if !c.enabled {
		return nil
	}
	return c.client.Del(ctx, key).Err()
}

func (c *CacheClient) Close() error {
	if !c.enabled || c.client == nil {
		return nil
	}
	return c.client.Close()
}

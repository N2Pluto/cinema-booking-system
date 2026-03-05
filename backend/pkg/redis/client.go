package redis

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

// Client wraps go-redis so the rest of the codebase imports only this pkg.
type Client struct {
	*redis.Client
}

// New creates a Redis client from the REDIS_URL environment variable.
// Expected format: redis://<host>:<port>  (no auth for dev; add :<password>@ for prod).
func New() (*Client, error) {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		url = "redis://localhost:6379"
	}

	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("redis: parse url: %w", err)
	}

	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis: ping: %w", err)
	}

	return &Client{rdb}, nil
}

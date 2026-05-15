package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

// NewClient returns a Redis client configured from address, password, and logical DB.
func NewClient(addr, password string, db int) *goredis.Client {
	return goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

// Ping verifies connectivity.
func Ping(ctx context.Context, c *goredis.Client) error {
	if err := c.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis: %w", err)
	}
	return nil
}

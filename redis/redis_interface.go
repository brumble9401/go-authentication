package redis

import (
	"context"
	"time"
)

type RedisClient interface {
    Ping(ctx context.Context) (string, error)
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    Get(ctx context.Context, key string) (string, error)
    Delete(ctx context.Context, key string) error
}
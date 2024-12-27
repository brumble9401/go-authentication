package session

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service struct {
    client *redis.Client
}

func NewService(client *redis.Client) *Service {
    return &Service{client: client}
}

func (s *Service) CreateSession(ctx context.Context, sessionID string, userID string, expiration time.Duration) error {
    return s.client.Set(ctx, sessionID, userID, expiration).Err()
}

func (s *Service) GetSession(ctx context.Context, sessionID string) (string, error) {
    return s.client.Get(ctx, sessionID).Result()
}

func (s *Service) DeleteSession(ctx context.Context, sessionID string) error {
    return s.client.Del(ctx, sessionID).Err()
}
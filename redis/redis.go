package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// NewRedisStore creates a new Redis session store
func NewRedisClient(redisAddr, redisPassword string, redisDB int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Ping Redis to check the connection
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Can not connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully!")
	return client
}
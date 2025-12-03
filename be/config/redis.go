package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(ctx context.Context) *redis.Client {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		log.Fatal("REDIS_URL not set")
	}

	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("failed parsing redis url: %v", err)
	}

	client := redis.NewClient(opt)

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	return client
}

package main

import (
	"context"
	"log/slog"
	"os"

	gowebly "github.com/gowebly/helpers"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func init() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file: %v", err)
		os.Exit(1)
	}

	ctx := context.Background()

	redisHost := gowebly.Getenv("REDIS_URL", "")
	opts, err := redis.ParseURL(redisHost)
	if err != nil {
		slog.Error("Error parsing URL to redis: %v", err)
		os.Exit(1)
	}

	redisClient = redis.NewClient(opts)
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		slog.Error("Error trying to ping redis: %v", err)
		os.Exit(1)
	}
}

func main() {
	if err := runServer(); err != nil {
		slog.Error("Failed to start server!", "details", err.Error())
		os.Exit(1)
	}
}

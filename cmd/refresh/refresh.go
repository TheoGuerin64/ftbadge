package main

import (
	"context"
	"ftbadge/internal/cache"
	"ftbadge/internal/ftapi"
	"ftbadge/internal/utils"
	"log"
)

func main() {
	redisURL := utils.MustGetEnv("REDIS_URL")

	redisStore, err := cache.NewRedisStore(redisURL)
	if err != nil {
		log.Fatalf("failed to setup Redis client: %v", err)
	}

	if _, err := ftapi.CacheAccessToken(context.Background(), redisStore); err != nil {
		log.Fatalf("failed to cache access token: %v", err)
	}
	log.Println("Access token cached successfully")
}

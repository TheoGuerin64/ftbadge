package main

import (
	"fmt"
	"ftbadge/internal/cache"
	"ftbadge/internal/ftcontext"
	"ftbadge/internal/ftvalidator"
	"ftbadge/internal/handlers"
	"ftbadge/internal/utils"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

const (
	defaultPort = "8080"
)

func setupRedisStore() (*cache.RedisStore, error) {
	redisURL := utils.MustGetEnv("REDIS_URL")

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse REDIS_URL %q: %w", redisURL, err)
	}
	client := redis.NewClient(opt)

	return cache.NewRedisStore(client), nil
}

func main() {
	port := utils.GetEnvWithDefault("PORT", defaultPort)

	redisStore, err := setupRedisStore()
	if err != nil {
		log.Fatalf("failed to setup Redis client: %v", err)
	}

	e := echo.New()
	e.Validator = ftvalidator.New()

	e.Use(ftcontext.Middleware(redisStore))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisableErrorHandler: true,
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())

	e.GET("/:login", handlers.ProfileHandler)

	log.Fatal(e.Start(":" + port))
}

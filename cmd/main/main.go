package main

import (
	"ftbadge/internal/cache"
	"ftbadge/internal/ftcontext"
	"ftbadge/internal/ftvalidator"
	"ftbadge/internal/handlers"
	"ftbadge/internal/utils"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	defaultPort = "8080"
)

func main() {
	port := utils.GetEnvWithDefault("PORT", defaultPort)
	redisURL := utils.MustGetEnv("REDIS_URL")

	redisStore, err := cache.NewRedisStore(redisURL)
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

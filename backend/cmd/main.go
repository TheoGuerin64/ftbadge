package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"ftbadge/internal/cache"
	"ftbadge/internal/ftvalidator"
	"ftbadge/internal/handlers"
)

func main() {
	redisClient, err := cache.NewRedisClient()
	if err != nil {
		log.Fatalf("failed to setup Redis client: %v", err)
	}

	e := echo.New()
	e.Validator = ftvalidator.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())

	e.GET("/:login", handlers.GetProfileHandler(redisClient))

	e.Logger.Fatal(e.Start(":3000"))
}

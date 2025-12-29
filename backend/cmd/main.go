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
	localClient, err := cache.NewLocalClient()
	if err != nil {
		log.Fatalf("failed to setup Redis client: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Validator = ftvalidator.New()

	e.Use(middleware.Recover())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Gzip())

	e.GET("/health", handlers.HealthCheckHandler)
	e.GET("/profile/:login", handlers.GetProfileHandler(localClient))

	e.Logger.Fatal(e.Start(":3000"))
}

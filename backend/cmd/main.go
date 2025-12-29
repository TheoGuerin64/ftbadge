package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"ftbadge/internal/cache"
	"ftbadge/internal/ftvalidator"
	"ftbadge/internal/handlers"
)

var (
	profileRateLimiterConfig = middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(10.0 / 60.0), Burst: 20, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			data := map[string]string{"error": "rate limit exceeded"}
			return context.JSON(http.StatusForbidden, data)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			data := map[string]string{"error": "rate limit exceeded"}
			return context.JSON(http.StatusTooManyRequests, data)
		},
	}
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
	e.GET("/profile/:login", handlers.GetProfileHandler(localClient), middleware.RateLimiterWithConfig(profileRateLimiterConfig))

	e.Logger.Fatal(e.Start(":3000"))
}

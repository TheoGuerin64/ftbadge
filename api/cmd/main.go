package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"

	"ftbadge/internal/cache"
	"ftbadge/internal/ftvalidator"
	"ftbadge/internal/handlers"
	"ftbadge/internal/utils"
)

func main() {
	port := utils.GetEnvWithDefault("PORT", "3000")
	sentryDSN := utils.MustGetEnv("SENTRY_DSN")

	sentryConfig := sentry.ClientOptions{
		Dsn: sentryDSN,
	}
	if err := sentry.Init(sentryConfig); err != nil {
		log.Fatalf("sentry initialization failed: %v", err)
	}

	localClient, err := cache.NewLocalClient()
	if err != nil {
		log.Fatalf("failed to setup Redis client: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Validator = ftvalidator.New()

	logger := zerolog.New(os.Stdout)
	requestLoggerConfig := middleware.RequestLoggerConfig{
		Skipper: func(ctx echo.Context) bool {
			return ctx.Request().URL.Path == "/health" && ctx.Response().Status == http.StatusOK
		},
		LogRequestID: true,
		LogMethod:    true,
		LogURI:       true,
		LogUserAgent: true,
		LogStatus:    true,
		LogError:     true,
		LogLatency:   true,
		LogValuesFunc: func(ctx echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Time("start_time", v.StartTime).
				Str("request_id", v.RequestID).
				Str("method", v.Method).
				Str("uri", v.URI).
				Str("user_agent", v.UserAgent).
				Int("status", v.Status).
				Err(v.Error).
				Dur("latency", v.Latency).
				Send()

			return nil
		},
	}
	sentryEchoConfig := sentryecho.Options{
		Repanic: true,
	}

	e.Use(middleware.RequestLoggerWithConfig(requestLoggerConfig))
	e.Use(middleware.Recover())
	e.Use(sentryecho.New(sentryEchoConfig))
	e.Use(middleware.Gzip())

	profileRateLimiterConfig := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(10.0 / 60.0), Burst: 20, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(ctx echo.Context, err error) error {
			data := map[string]string{"error": "rate limit exceeded"}
			return ctx.JSON(http.StatusForbidden, data)
		},
		DenyHandler: func(ctx echo.Context, identifier string, err error) error {
			data := map[string]string{"error": "rate limit exceeded"}
			return ctx.JSON(http.StatusTooManyRequests, data)
		},
	}

	e.GET("/health", handlers.HealthCheckHandler)
	e.GET("/profile/:login", handlers.GetProfileHandler(localClient), middleware.RateLimiterWithConfig(profileRateLimiterConfig))

	e.Logger.Fatal(e.Start(":" + port))
}

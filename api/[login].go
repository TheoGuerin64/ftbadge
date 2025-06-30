package handler

import (
	"ftbadge/server/cache"
	"ftbadge/server/ftvalidator"
	"ftbadge/server/handlers"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	redisStore, err := cache.NewRedisClient()
	if err != nil {
		log.Fatalf("failed to setup Redis client: %v", err)
	}

	e := echo.New()
	e.Validator = ftvalidator.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())

	e.GET("/:login", handlers.GetProfileHandler(redisStore))

	e.ServeHTTP(w, r)
}

package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealthCheckHandler(ctx echo.Context) error {
	ctx.Response().Header().Set("Cache-Control", "no-store, no-cache, max-age=0")
	data := map[string]string{"status": "ok"}
	return ctx.JSON(http.StatusOK, data)
}

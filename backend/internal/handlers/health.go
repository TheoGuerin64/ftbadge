package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealthCheckHandler(ctx echo.Context) error {
	data := map[string]string{"status": "ok"}
	return ctx.JSON(http.StatusOK, data)
}

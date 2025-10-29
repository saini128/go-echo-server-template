package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealthRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "You have reached the Center of Excellence TIET-UQ API server."})
	})

}

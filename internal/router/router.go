package router

import (
	"github.com/labstack/echo/v4"
	"github.com/wando-world/wando-sso/internal/api"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/", api.HomeHandler)
}

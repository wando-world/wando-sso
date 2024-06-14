package router

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	apiGroup := e.Group("/sso/api/v1")

	SetupEnvRoutes(apiGroup.Group("/env"))
	SetupUserRoutes(apiGroup.Group("/user"))
}

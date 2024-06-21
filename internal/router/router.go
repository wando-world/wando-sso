package router

import (
	"github.com/labstack/echo/v4"
	"github.com/wando-world/wando-sso/internal/router/middleware"
	"github.com/wando-world/wando-sso/internal/utils"
)

func SetupRoutes(e *echo.Echo, jwtUtils utils.IJwt) {
	apiGroup := e.Group("/sso/api/v1")

	SetupEnvRoutes(apiGroup.Group("/env"))
	SetupAuthRoutes(apiGroup.Group("/auth"), jwtUtils)

	userGroup := apiGroup.Group("/user")
	userGroup.Use(middleware.JwtMiddleware(jwtUtils))
	SetupUserRoutes(userGroup)
}

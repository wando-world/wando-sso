package router

import (
	"github.com/labstack/echo/v4"
	"github.com/wando-world/wando-sso/internal/api"
)

func SetupUserRoutes(g *echo.Group) {
	g.POST("", api.CreateUserHandler)
}

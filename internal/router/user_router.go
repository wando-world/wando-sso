package router

import (
	"github.com/labstack/echo/v4"
	"github.com/wando-world/wando-sso/internal/api/wire"
)

func SetupUserRoutes(g *echo.Group) {
	userHandler := wire.InitUserHandler()

	g.POST("", userHandler.CreateUser)
}

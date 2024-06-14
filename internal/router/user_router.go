package router

import (
	"github.com/labstack/echo/v4"
	"github.com/wando-world/wando-sso/internal/api"
	"github.com/wando-world/wando-sso/internal/db"
	"github.com/wando-world/wando-sso/internal/mappers"
)

func SetupUserRoutes(g *echo.Group) {
	userMapper := mappers.NewUserMapper()
	userRepository := db.NewUserRepository(db.DB)
	userHandler := api.NewUserHandler(userMapper, userRepository)

	g.POST("", userHandler.CreateUser)
}

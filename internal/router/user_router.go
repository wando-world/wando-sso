package router

import (
	"github.com/labstack/echo/v4"
	"github.com/wando-world/wando-sso/internal/api"
	"github.com/wando-world/wando-sso/internal/db"
	"github.com/wando-world/wando-sso/internal/mappers"
	"github.com/wando-world/wando-sso/internal/utils"
)

func SetupUserRoutes(g *echo.Group) {
	passwordUtils := utils.NewPasswordUtils()
	userMapper := mappers.NewUserMapper(passwordUtils)
	userRepository := db.NewUserRepository(db.DB)
	userHandler := api.NewUserHandler(userMapper, userRepository)

	g.POST("", userHandler.CreateUser)
}

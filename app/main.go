package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wando-world/wando-sso/auth"
	"github.com/wando-world/wando-sso/internal/config"
	"github.com/wando-world/wando-sso/internal/repository/postgresql"
	"github.com/wando-world/wando-sso/internal/rest"
	restmiddleware "github.com/wando-world/wando-sso/internal/rest/middleware"
	"github.com/wando-world/wando-sso/user"
)

func main() {
	// prepare config
	cfg := config.New()

	// prepare postgresql database
	postgresql.InitDB(cfg.DbUrl)
	defer postgresql.CloseDB()

	// prepare echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = restmiddleware.NewValidator()

	// prepare api group
	apiGroup := e.Group("/sso/api/v1")

	// prepare auth api
	authGroup := apiGroup.Group("/auth")
	authGroup.Use(restmiddleware.RtkMiddleware())

	// prepare user api
	userGroup := apiGroup.Group("/user")
	userGroup.Use(restmiddleware.AtkMiddleware())

	// prepare repository
	authRepo := postgresql.NewAuthRepository(postgresql.DB)
	userRepo := postgresql.NewUserRepository(postgresql.DB)

	// build service layer
	authSvc := auth.NewService(authRepo)
	userSvc := user.NewService(userRepo)

	rest.NewAuthHandler(authGroup, authSvc)
	rest.NewUserHandler(userGroup, userSvc)

	e.Logger.Fatal(e.Start(cfg.Port))
}

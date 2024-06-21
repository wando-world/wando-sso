package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/wando-world/wando-sso/internal/utils"
	"net/http"
)

func JwtMiddleware(jwtUtils utils.IJwt) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusBadRequest, "인증된 유저가 아니군요!\n로그인이 필요해요!")
		},
		SigningKey: jwtUtils.(*utils.JwtUtils).AtkSecret,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(utils.Claims)
		},
	})
}

package middleware

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/wando-world/wando-sso/internal/config"
	"github.com/wando-world/wando-sso/utils"
	"net/http"
	"strings"
)

func AtkMiddleware() echo.MiddlewareFunc {
	jwtUtils := utils.NewJwtUtils(config.GetConfig().ATKSecret, config.GetConfig().RTKSecret)

	return echojwt.WithConfig(echojwt.Config{
		ErrorHandler: func(c echo.Context, err error) error {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return echo.NewHTTPError(http.StatusUnauthorized, "토큰이 만료되었습니다.\nFE 에서 토큰을 리프레쉬를 해주세요.")
			}
			return echo.NewHTTPError(http.StatusBadRequest, "인증된 유저가 아니군요!\n로그인이 필요해요!")
		},
		SigningKey:    jwtUtils.GetAtkSecret(),
		SigningMethod: "HS512",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(utils.Claims)
		},
		Skipper: func(c echo.Context) bool {
			if strings.HasPrefix(c.Path(), "/sso/api/v1/user") && c.Request().Method == http.MethodPost { // 회원 가입만 skip
				return true
			}
			return false
		},
	})
}

func RtkMiddleware() echo.MiddlewareFunc {
	jwtUtils := utils.NewJwtUtils(config.GetConfig().ATKSecret, config.GetConfig().RTKSecret)
	return echojwt.WithConfig(echojwt.Config{
		ErrorHandler: func(c echo.Context, err error) error {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return echo.NewHTTPError(http.StatusUnauthorized, "인증 시간이 다 되어서 로그인이 필요해요!")
			}
			return echo.NewHTTPError(http.StatusBadRequest, "인증된 유저가 아니군요!\n로그인이 필요해요!")
		},
		SigningKey:    jwtUtils.GetRtkSecret(),
		SigningMethod: "HS512",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(utils.Claims)
		},
		Skipper: func(c echo.Context) bool {
			if strings.HasSuffix(c.Path(), "refresh") {
				return false
			}
			return true
		},
	})
}

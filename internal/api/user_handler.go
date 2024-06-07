package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/wando-world/wando-sso/internal/api/models"
	"github.com/wando-world/wando-sso/internal/db"
	"github.com/wando-world/wando-sso/internal/mappers"
	"gorm.io/gorm"
	"net/http"
)

func CreateUserHandler(c echo.Context) error {
	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "입력값을 확인해주세요.")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	mapper := mappers.NewUserMapper()
	userRepository := db.NewUserRepository(db.DB)

	user, err := mapper.CreateUserRequestToUser(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "암호화 중 에러가 발생했습니다!\n잠시뒤 진행해 주세요!")
	}

	err = userRepository.CreateUser(&user)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return echo.NewHTTPError(http.StatusConflict, "id 가 이미 있습니다!")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "가입 실패!")
	}
	return c.JSON(http.StatusCreated, user.Nickname)
}

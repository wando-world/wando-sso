package rest

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/wando-world/wando-sso/domain"
	"github.com/wando-world/wando-sso/user"
	"github.com/wando-world/wando-sso/utils"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type UserService interface {
	SignupUser(ctx context.Context, req user.SignupUserReq) error
	FindSelfById(ctx context.Context, id uint) (*domain.User, error)
}

type UserHandler struct {
	UserService UserService
}

func NewUserHandler(g *echo.Group, us UserService) {
	handler := &UserHandler{
		UserService: us,
	}
	g.POST("", handler.SignupUser)
	g.GET("", handler.FindSelfById)
}

func (uh *UserHandler) SignupUser(c echo.Context) error {
	var req user.SignupUserReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "입력값을 확인해주세요.")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 1*time.Second)
	defer cancel()

	err := uh.UserService.SignupUser(ctx, req)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return echo.NewHTTPError(http.StatusConflict, "id 가 이미 있습니다!")
	} else if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return echo.NewHTTPError(http.StatusRequestTimeout, "요청 시간이 초과되었습니다..\n나중에 다시 시도해 주세요..")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "가입 실패!")
	}
	return c.JSON(http.StatusCreated, nil)
}

func (uh *UserHandler) FindSelfById(c echo.Context) error {
	loginUser := c.Get("user").(*jwt.Token)
	claims := loginUser.Claims.(*utils.Claims)
	id := claims.Id

	ctx, cancel := context.WithTimeout(c.Request().Context(), 1*time.Second)
	defer cancel()

	foundUser, err := uh.UserService.FindSelfById(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "유저가 없습니다.")
	} else if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return echo.NewHTTPError(http.StatusRequestTimeout, "요청 시간이 초과되었습니다..\n나중에 다시 시도해 주세요..")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "서버가 문제가 있어요.\n어떻게 하셨을때 에러가 났는지 문의에 남겨주세요!")
	}

	return c.JSON(http.StatusOK, user.FindSelfResponse{
		Nickname: foundUser.Nickname,
		UserID:   foundUser.UserID,
		Email:    foundUser.Email,
	})
}

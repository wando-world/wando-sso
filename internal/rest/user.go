package rest

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/wando-world/wando-sso/domain"
	apiModel "github.com/wando-world/wando-sso/internal/rest/dto"
	"github.com/wando-world/wando-sso/internal/rest/mappers"
	"github.com/wando-world/wando-sso/internal/rest/middleware"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type UserService interface {
	SignupUser(ctx context.Context, u *domain.User) error
	FindSelfById(ctx context.Context, id uint) (*domain.User, error)
}

type UserHandler struct {
	UserService UserService
	UserMapper  mappers.IUserMapper
}

func NewUserHandler(g *echo.Group, us UserService, m mappers.IUserMapper) {
	handler := &UserHandler{
		UserService: us,
		UserMapper:  m,
	}
	g.POST("", handler.SignupUser)
	g.GET("", handler.FindSelfById)
}

func (uh *UserHandler) SignupUser(c echo.Context) error {
	var req apiModel.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "입력값을 확인해주세요.")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := uh.UserMapper.SignupUserRequestToUser(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "암호화 중 에러가 발생했습니다!\n잠시뒤 진행해 주세요!")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 1*time.Second)
	defer cancel()

	err = uh.UserService.SignupUser(ctx, &user)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return echo.NewHTTPError(http.StatusConflict, "id 가 이미 있습니다!")
	} else if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return echo.NewHTTPError(http.StatusRequestTimeout, "요청 시간이 초과되었습니다..\n나중에 다시 시도해 주세요..")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "가입 실패!")
	}
	return c.JSON(http.StatusCreated, user.Nickname)
}

func (uh *UserHandler) FindSelfById(c echo.Context) error {
	loginUser := c.Get("user").(*jwt.Token)
	claims := loginUser.Claims.(*middleware.Claims)
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

	return c.JSON(http.StatusOK, apiModel.FindSelfResponse{
		Nickname: foundUser.Nickname,
		UserID:   foundUser.UserID,
		Email:    foundUser.Email,
	})
}

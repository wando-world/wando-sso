package api

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	apiModel "github.com/wando-world/wando-sso/internal/api/models"
	"github.com/wando-world/wando-sso/internal/db"
	"github.com/wando-world/wando-sso/internal/mappers"
	"github.com/wando-world/wando-sso/internal/models"
	"github.com/wando-world/wando-sso/internal/utils"
	"gorm.io/gorm"
	"net/http"
)

type IUserHandler interface {
	SignupUser(c echo.Context) error
	FindSelfById(c echo.Context) error
}

// UserHandler 구조체에 의존성을 저장
type UserHandler struct {
	Mapper         mappers.IUserMapper
	UserRepository db.IUserRepository
}

// NewUserHandler 생성자 함수는 의존성을 주입
func NewUserHandler(mapper mappers.IUserMapper, userRepository db.IUserRepository) IUserHandler {
	return &UserHandler{
		Mapper:         mapper,
		UserRepository: userRepository,
	}
}

// SignupUser 회원가입
func (h *UserHandler) SignupUser(c echo.Context) error {
	var req apiModel.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "입력값을 확인해주세요.")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.Mapper.CreateUserRequestToUser(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "암호화 중 에러가 발생했습니다!\n잠시뒤 진행해 주세요!")
	}

	err = h.UserRepository.CreateUser(&user)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return echo.NewHTTPError(http.StatusConflict, "id 가 이미 있습니다!")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "가입 실패!")
	}
	return c.JSON(http.StatusCreated, user.Nickname)
}

func (h *UserHandler) FindSelfById(c echo.Context) error {
	loginUser := c.Get("user").(*jwt.Token)
	claims := loginUser.Claims.(*utils.Claims)
	id := claims.Id

	user := models.User{
		Model: gorm.Model{ID: id},
	}

	foundUser, err := h.UserRepository.FindUserByID(&user)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "유저가 없습니다.")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "서버가 문제가 있어요.\n어떻게 하셨을때 에러가 났는지 문의에 남겨주세요!")
	}

	return c.JSON(http.StatusOK, apiModel.FindSelfResponse{
		Nickname: foundUser.Nickname,
		UserID:   foundUser.UserID,
		Email:    foundUser.Email,
	})
}

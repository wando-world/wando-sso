package api

import (
	"encoding/base64"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/wando-world/wando-sso/internal/api/models"
	"github.com/wando-world/wando-sso/internal/db"
	"github.com/wando-world/wando-sso/internal/mappers"
	"github.com/wando-world/wando-sso/internal/utils"
	"gorm.io/gorm"
	"net/http"
)

type IAuthHandler interface {
	Login(c echo.Context) error
}

type AuthHandler struct {
	PasswordUtils  utils.IPasswordUtils
	JwtUtils       utils.IJwt
	Mapper         mappers.IAuthMapper
	UserRepository db.IUserRepository
}

func NewAuthHandler(passwordUtils utils.IPasswordUtils, jwtUtils utils.IJwt, mapper mappers.IAuthMapper, userRepository db.IUserRepository) IAuthHandler {
	return &AuthHandler{
		PasswordUtils:  passwordUtils,
		JwtUtils:       jwtUtils,
		Mapper:         mapper,
		UserRepository: userRepository,
	}
}

func (a AuthHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "입력값을 확인해주세요.")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := a.Mapper.LoginRequestToUser(req)

	foundUser, err := a.UserRepository.FindUserForLogin(&user)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "유저가 없습니다.")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "서버가 문제가 있어요.\n어떻게 하셨을때 에러가 났는지 문의에 남겨주세요!")
	}

	// 비밀번호 체크
	decoded, err := base64.RawStdEncoding.DecodeString(foundUser.Salt)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "서버가 문제가 있어요.\n어떻게 하셨을때 에러가 났는지 문의에 남겨주세요!")
	}
	if ok := a.PasswordUtils.VerifyPassword(req.Password, foundUser.Password, decoded); !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "아이디 또는 비밀번호를 확인해주세요!")
	}

	atk, err := a.JwtUtils.GenerateATK(foundUser.ID, foundUser.Role)
	if err != nil {
		log.Errorf("atk 발급 에러: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "로그인중 서버 에러가 발생했어요ㅠㅠ\n문의를 남겨주세요!")
	}

	rtk, err := a.JwtUtils.GenerateRTK(foundUser.ID)
	if err != nil {
		log.Errorf("rtk 발급 에러: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "로그인중 서버 에러가 발생했어요ㅠㅠ\n문의를 남겨주세요!")
	}

	return c.JSON(http.StatusOK, models.LoginResponse{
		ATK: atk,
		RTK: rtk,
	})
}

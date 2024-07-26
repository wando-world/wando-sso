package rest

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/wando-world/wando-sso/domain"
	apiModel "github.com/wando-world/wando-sso/internal/rest/dto"
	"github.com/wando-world/wando-sso/internal/rest/mappers"
	"github.com/wando-world/wando-sso/internal/rest/middleware"
	"github.com/wando-world/wando-sso/utils"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type AuthService interface {
	Login(ctx context.Context, userId, verifiedCode string) (*domain.User, error)
	RefreshAtk(ctx context.Context, id uint) (*domain.User, error)
}

type AuthHandler struct {
	AuthService   AuthService
	AuthMapper    mappers.IAuthMapper
	PasswordUtils utils.IPasswordUtils
	JwtService    middleware.IJwt
}

func NewAuthHandler(g *echo.Group, as AuthService, m mappers.IAuthMapper, p utils.IPasswordUtils, js middleware.IJwt) {
	handler := &AuthHandler{
		AuthService:   as,
		AuthMapper:    m,
		PasswordUtils: p,
		JwtService:    js,
	}
	g.POST("/login", handler.Login)
}

func (ah *AuthHandler) Login(c echo.Context) error {
	var req apiModel.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "입력값을 확인해주세요.")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := ah.AuthMapper.LoginRequestToUser(req)

	ctx, cancel := context.WithTimeout(c.Request().Context(), 1*time.Second)
	defer cancel()

	foundUser, err := ah.AuthService.Login(ctx, user.UserID, user.VerifiedCode)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "유저가 없습니다.")
	} else if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return echo.NewHTTPError(http.StatusRequestTimeout, "요청 시간이 초과되었습니다..\n나중에 다시 시도해 주세요..")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "서버가 문제가 있어요.\n어떻게 하셨을때 에러가 났는지 문의에 남겨주세요!")
	}

	// 비밀번호 체크
	decoded, err := base64.RawStdEncoding.DecodeString(foundUser.Salt)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "서버가 문제가 있어요.\n어떻게 하셨을때 에러가 났는지 문의에 남겨주세요!")
	}
	if ok := ah.PasswordUtils.VerifyPassword(req.Password, foundUser.Password, decoded); !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "아이디 또는 비밀번호를 확인해주세요!")
	}

	// ATK 발급
	atk, err := ah.JwtService.GenerateATK(foundUser.ID, foundUser.Role)
	if err != nil {
		log.Errorf("atk 발급 에러: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "로그인중 서버 에러가 발생했어요ㅠㅠ\n문의를 남겨주세요!")
	}

	// RTK 발급
	rtk, err := ah.JwtService.GenerateRTK(foundUser.ID)
	if err != nil {
		log.Errorf("rtk 발급 에러: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "로그인중 서버 에러가 발생했어요ㅠㅠ\n문의를 남겨주세요!")
	}

	return c.JSON(http.StatusOK, apiModel.LoginResponse{
		ATK: atk,
		RTK: rtk,
	})
}

func (ah *AuthHandler) RefreshAtk(c echo.Context) error {
	loginUser := c.Get("user").(*jwt.Token)
	claims := loginUser.Claims.(*middleware.Claims)
	id := claims.Id

	ctx, cancel := context.WithTimeout(c.Request().Context(), 1*time.Second)
	defer cancel()

	foundUser, err := ah.AuthService.RefreshAtk(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "유저가 없습니다.")
	} else if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return echo.NewHTTPError(http.StatusRequestTimeout, "요청 시간이 초과되었습니다..\n나중에 다시 시도해 주세요..")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "서버가 문제가 있어요.\n어떻게 하셨을때 에러가 났는지 문의에 남겨주세요!")
	}

	atk, err := ah.JwtService.GenerateATK(foundUser.ID, foundUser.Role)
	if err != nil {
		log.Errorf("atk 발급 에러: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "인증 갱신 중 서버 에러가 발생했어요ㅠㅠ\n문의를 남겨주세요!")
	}
	return c.JSON(http.StatusOK, apiModel.RefreshAtkResponse{
		ATK: atk,
	})
}

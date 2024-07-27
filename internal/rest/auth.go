package rest

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/wando-world/wando-sso/auth"
	"github.com/wando-world/wando-sso/domain"
	"github.com/wando-world/wando-sso/internal/rest/dto"
	"github.com/wando-world/wando-sso/utils"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type AuthService interface {
	FindUser(ctx context.Context, req auth.FindUserReq) (*domain.User, error)
	CheckPassword(salt, encryptPassword, plaintextPassword string) (bool, error)
	GenerateJWT(id uint, role string) (*auth.GenerateJWTRes, error)
	RefreshAtk(ctx context.Context, id uint) (*auth.RefreshAtkRes, error)
}

type AuthHandler struct {
	AuthService AuthService
}

func NewAuthHandler(g *echo.Group, as AuthService) {
	handler := &AuthHandler{
		AuthService: as,
	}
	g.POST("/login", handler.Login)
}

func (ah *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "입력값을 확인해주세요.")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 1*time.Second)
	defer cancel()

	foundUser, err := ah.AuthService.FindUser(ctx, auth.FindUserReq{
		UserId:       req.UserID,
		VerifiedCode: req.VerifiedCode,
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "유저가 없습니다.")
	} else if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return echo.NewHTTPError(http.StatusRequestTimeout, "요청 시간이 초과되었습니다..\n나중에 다시 시도해 주세요..")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "서버가 문제가 있어요.\n어떻게 하셨을때 에러가 났는지 문의에 남겨주세요!")
	}

	// 비밀번호 체크
	checkPassword, err := ah.AuthService.CheckPassword(foundUser.Salt, foundUser.Password, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "서버가 문제가 있어요.\n어떻게 하셨을때 에러가 났는지 문의에 남겨주세요!")
	}
	if !checkPassword {
		return echo.NewHTTPError(http.StatusBadRequest, "아이디 또는 비밀번호를 확인해주세요!")
	}

	// JWT 발급
	jwts, err := ah.AuthService.GenerateJWT(foundUser.ID, foundUser.Role)
	if err != nil {
		log.Errorf("jwt 발급 에러: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "로그인중 서버 에러가 발생했어요ㅠㅠ\n문의를 남겨주세요!")
	}

	return c.JSON(http.StatusOK, jwts)
}

func (ah *AuthHandler) RefreshAtk(c echo.Context) error {
	loginUser := c.Get("user").(*jwt.Token)
	claims := loginUser.Claims.(*utils.Claims)
	id := claims.Id

	ctx, cancel := context.WithTimeout(c.Request().Context(), 1*time.Second)
	defer cancel()

	atk, err := ah.AuthService.RefreshAtk(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "유저가 없습니다.")
	} else if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return echo.NewHTTPError(http.StatusRequestTimeout, "요청 시간이 초과되었습니다..\n나중에 다시 시도해 주세요..")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "서버가 문제가 있어요.\n어떻게 하셨을때 에러가 났는지 문의에 남겨주세요!")
	}

	return c.JSON(http.StatusOK, atk)
}

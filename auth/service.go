package auth

import (
	"context"
	"encoding/base64"
	"github.com/wando-world/wando-sso/domain"
	"github.com/wando-world/wando-sso/internal/config"
	"github.com/wando-world/wando-sso/utils"
)

type AuthRepository interface {
	FindUserForLogin(ctx context.Context, userId, verifiedCode string) (*domain.User, error)
	FindUserById(ctx context.Context, id uint) (*domain.User, error)
}

type Service struct {
	authRepo AuthRepository
}

func NewService(ar AuthRepository) *Service {
	return &Service{
		authRepo: ar,
	}
}

func (s *Service) FindUser(ctx context.Context, req FindUserReq) (*domain.User, error) {
	return s.authRepo.FindUserForLogin(ctx, req.UserId, req.VerifiedCode)
}

func (s *Service) CheckPassword(salt, encryptPassword, plaintextPassword string) (bool, error) {
	passwordUtils := utils.NewPasswordUtils()
	decoded, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return false, err
	}
	if ok := passwordUtils.VerifyPassword(plaintextPassword, encryptPassword, decoded); !ok {
		return false, nil
	}

	return true, nil
}

func (s *Service) GenerateJWT(id uint, role string) (*GenerateJWTRes, error) {
	atk, err := generateATK(id, role)
	if err != nil {
		return nil, err
	}

	rtk, err := generateRTK(id)
	if err != nil {
		return nil, err
	}

	return &GenerateJWTRes{
		ATK: atk,
		RTK: rtk,
	}, nil
}

func (s *Service) RefreshAtk(ctx context.Context, id uint) (*RefreshAtkRes, error) {
	found, err := s.authRepo.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	atk, err := generateATK(found.ID, found.Role)
	if err != nil {
		return nil, err
	}

	return &RefreshAtkRes{
		ATK: atk,
	}, nil
}

func generateATK(id uint, role string) (string, error) {
	jwtUtils := utils.NewJwtUtils(config.GetConfig().ATKSecret, config.GetConfig().RTKSecret)
	atk, err := jwtUtils.GenerateATK(id, role)
	if err != nil {
		return "", err
	}
	return atk, nil
}

func generateRTK(id uint) (string, error) {
	jwtUtils := utils.NewJwtUtils(config.GetConfig().ATKSecret, config.GetConfig().RTKSecret)
	rtk, err := jwtUtils.GenerateRTK(id)
	if err != nil {
		return "", err
	}
	return rtk, nil
}

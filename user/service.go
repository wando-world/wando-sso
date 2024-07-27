package user

import (
	"context"
	"encoding/base64"
	"github.com/wando-world/wando-sso/domain"
	"github.com/wando-world/wando-sso/utils"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u *domain.User) error
	FindUserById(ctx context.Context, id uint) (*domain.User, error)
}

type Service struct {
	userRepo UserRepository
}

func NewService(ur UserRepository) *Service {
	return &Service{
		userRepo: ur,
	}
}

func (s *Service) SignupUser(ctx context.Context, req SignupUserReq) error {
	user, err := signupUserRequestToUser(req)
	if err != nil {
		return err
	}

	return s.userRepo.CreateUser(ctx, &user)
}

func (s *Service) FindSelfById(ctx context.Context, id uint) (foundUser *domain.User, err error) {
	foundUser, err = s.userRepo.FindUserById(ctx, id)
	return
}

func signupUserRequestToUser(req SignupUserReq) (domain.User, error) {
	passwordUtils := utils.NewPasswordUtils()
	user := domain.User{
		Nickname:     req.Nickname,
		UserID:       req.UserID,
		Email:        req.Email,
		VerifiedCode: req.VerifiedCode,
		Role:         "GENERAL",
		Password:     req.Password,
	}

	salt, err := passwordUtils.GenerateSalt()
	if err != nil {
		return domain.User{}, err
	}

	user.Password = passwordUtils.HashPassword(user.Password, salt)
	user.Salt = base64.RawStdEncoding.EncodeToString(salt)

	return user, nil
}

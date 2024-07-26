package user

import (
	"context"
	"github.com/wando-world/wando-sso/domain"
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

func (s *Service) SignupUser(ctx context.Context, u *domain.User) error {
	return s.userRepo.CreateUser(ctx, u)
}

func (s *Service) FindSelfById(ctx context.Context, id uint) (foundUser *domain.User, err error) {
	foundUser, err = s.userRepo.FindUserById(ctx, id)
	return
}

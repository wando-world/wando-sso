package auth

import (
	"context"
	"github.com/wando-world/wando-sso/domain"
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

func (s *Service) Login(ctx context.Context, userId, verifiedCode string) (*domain.User, error) {
	return s.authRepo.FindUserForLogin(ctx, userId, verifiedCode)
}

func (s *Service) RefreshAtk(ctx context.Context, id uint) (*domain.User, error) {
	return s.authRepo.FindUserById(ctx, id)
}

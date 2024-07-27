package postgresql

import (
	"context"
	"github.com/wando-world/wando-sso/domain"
	"gorm.io/gorm"
	"sync"
	"testing"
)

type AuthRepository struct {
	DB *gorm.DB
}

var (
	authInstance *AuthRepository
	authOnce     sync.Once
)

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	if testing.Testing() {
		// 테스트 환경에서는 새 인스턴스 반환
		return &AuthRepository{DB: db}
	}

	authOnce.Do(func() {
		authInstance = &AuthRepository{DB: db}
	})
	return authInstance
}

func (ar *AuthRepository) FindUserForLogin(ctx context.Context, userId, verifiedCode string) (*domain.User, error) {
	var resultUser domain.User
	if err := ar.DB.WithContext(ctx).Where(&domain.User{UserID: userId, VerifiedCode: verifiedCode}).First(&resultUser).Error; err != nil {
		return nil, err
	}
	return &resultUser, nil
}

func (ar *AuthRepository) FindUserById(ctx context.Context, id uint) (*domain.User, error) {
	var resultUser domain.User
	result := ar.DB.WithContext(ctx).First(&resultUser, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &resultUser, nil
}

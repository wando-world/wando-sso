package postgresql

import (
	"context"
	"github.com/wando-world/wando-sso/domain"
	"gorm.io/gorm"
	"sync"
	"testing"
)

type UserRepository struct {
	DB *gorm.DB
}

var (
	userInstance *UserRepository
	userOnce     sync.Once
)

func NewUserRepository(db *gorm.DB) *UserRepository {
	if testing.Testing() {
		// 테스트 환경에서는 새 인스턴스 반환
		return &UserRepository{DB: db}
	}
	userOnce.Do(func() {
		userInstance = &UserRepository{DB: db}
	})
	return userInstance
}

func (ur *UserRepository) CreateUser(ctx context.Context, u *domain.User) error {
	return ur.DB.WithContext(ctx).Create(u).Error
}

func (ur *UserRepository) FindUserById(ctx context.Context, id uint) (*domain.User, error) {
	var resultUser domain.User
	result := ur.DB.WithContext(ctx).First(&resultUser, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &resultUser, nil
}

package postgresql

import (
	"context"
	"github.com/wando-world/wando-sso/domain"
	"gorm.io/gorm"
	"sync"
)

type UserRepository struct {
	DB *gorm.DB
}

var (
	userInstance *UserRepository
	userOnce     sync.Once
)

func NewUserRepository(db *gorm.DB) *UserRepository {
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

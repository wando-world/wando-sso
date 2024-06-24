package db

import (
	"github.com/wando-world/wando-sso/internal/models"
	"gorm.io/gorm"
	"sync"
)

type IUserRepository interface {
	CreateUser(user *models.User) error
	FindUserForLogin(user *models.User) (*models.User, error)
	FindUserByID(user *models.User) (*models.User, error)
}

type UserRepository struct {
	DB *gorm.DB
}

var (
	instance *UserRepository
	once     sync.Once
)

func NewUserRepository(db *gorm.DB) *UserRepository {
	once.Do(func() {
		instance = &UserRepository{
			DB: db,
		}
	})
	return instance
}

func (repo *UserRepository) CreateUser(user *models.User) error {
	return repo.DB.Create(user).Error
}

func (repo *UserRepository) FindUserForLogin(user *models.User) (*models.User, error) {
	var resultUser models.User
	if err := repo.DB.Where(&models.User{UserID: user.UserID, VerifiedCode: user.VerifiedCode}).First(&resultUser).Error; err != nil {
		return nil, err
	}
	return &resultUser, nil
}

func (repo *UserRepository) FindUserByID(user *models.User) (*models.User, error) {
	var resultUser models.User
	result := repo.DB.First(&resultUser, user.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &resultUser, nil
}

func (repo *UserRepository) UpdateUser(user *models.User) error {
	return repo.DB.Save(user).Error
}

func (repo *UserRepository) DeleteUser(id uint) error {
	result := repo.DB.Delete(&models.User{}, id)
	return result.Error
}

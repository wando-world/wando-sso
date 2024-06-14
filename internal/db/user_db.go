package db

import (
	"github.com/wando-world/wando-sso/internal/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user *models.User) error
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) CreateUser(user *models.User) error {
	return repo.DB.Create(user).Error
}

func (repo *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := repo.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) UpdateUser(user *models.User) error {
	return repo.DB.Save(user).Error
}

func (repo *UserRepository) DeleteUser(id uint) error {
	result := repo.DB.Delete(&models.User{}, id)
	return result.Error
}

package mappers

import (
	"encoding/base64"
	apiModels "github.com/wando-world/wando-sso/internal/api/models"
	"github.com/wando-world/wando-sso/internal/models"
	"github.com/wando-world/wando-sso/internal/utils"
)

type IUserMapper interface {
	CreateUserRequestToUser(req apiModels.CreateUserRequest) (models.User, error)
}

type UserMapper struct {
	passwordUtils utils.IPasswordUtils
}

func NewUserMapper(p utils.IPasswordUtils) *UserMapper {
	return &UserMapper{passwordUtils: p}
}

func (m *UserMapper) CreateUserRequestToUser(req apiModels.CreateUserRequest) (models.User, error) {
	user := models.User{
		Nickname:     req.Nickname,
		UserID:       req.UserID,
		Email:        req.Email,
		VerifiedCode: req.VerifiedCode,
		Role:         "GENERAL",
		Password:     req.Password,
	}

	salt, err := m.passwordUtils.GenerateSalt()
	if err != nil {
		return models.User{}, err
	}

	user.Password = m.passwordUtils.HashPassword(user.Password, salt)
	user.Salt = base64.RawStdEncoding.EncodeToString(salt)

	return user, nil
}

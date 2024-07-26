package mappers

import (
	"encoding/base64"
	"github.com/wando-world/wando-sso/domain"
	apiModels "github.com/wando-world/wando-sso/internal/rest/dto"
	"github.com/wando-world/wando-sso/utils"
)

type IUserMapper interface {
	SignupUserRequestToUser(req apiModels.CreateUserRequest) (domain.User, error)
}

type UserMapper struct {
	passwordUtils utils.IPasswordUtils
}

func NewUserMapper(p utils.IPasswordUtils) *UserMapper {
	return &UserMapper{passwordUtils: p}
}

func (m *UserMapper) SignupUserRequestToUser(req apiModels.CreateUserRequest) (domain.User, error) {
	user := domain.User{
		Nickname:     req.Nickname,
		UserID:       req.UserID,
		Email:        req.Email,
		VerifiedCode: req.VerifiedCode,
		Role:         "GENERAL",
		Password:     req.Password,
	}

	salt, err := m.passwordUtils.GenerateSalt()
	if err != nil {
		return domain.User{}, err
	}

	user.Password = m.passwordUtils.HashPassword(user.Password, salt)
	user.Salt = base64.RawStdEncoding.EncodeToString(salt)

	return user, nil
}

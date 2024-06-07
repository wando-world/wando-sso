package mappers

import (
	"encoding/base64"
	apiModels "github.com/wando-world/wando-sso/internal/api/models"
	"github.com/wando-world/wando-sso/internal/models"
	"github.com/wando-world/wando-sso/internal/utils"
)

type DefaultUserMapper struct{}

func NewUserMapper() IUserMapper {
	return &DefaultUserMapper{}
}

func (m *DefaultUserMapper) CreateUserRequestToUser(req apiModels.CreateUserRequest) (models.User, error) {
	user := models.User{
		Nickname:     req.Nickname,
		UserID:       req.UserID,
		Email:        req.Email,
		VerifiedCode: req.VerifiedCode,
		Role:         "GENERAL",
	}

	salt, err := utils.GenerateSalt()
	if err != nil {
		return models.User{}, err
	}

	user.Password = utils.HashPassword(user.Password, salt)
	user.Salt = base64.RawStdEncoding.EncodeToString(salt)

	return user, nil
}

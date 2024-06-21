package mappers

import (
	apiModels "github.com/wando-world/wando-sso/internal/api/models"
	"github.com/wando-world/wando-sso/internal/models"
)

type IAuthMapper interface {
	LoginRequestToUser(req apiModels.LoginRequest) models.User
}

type AuthMapper struct {
}

func NewAuthMapper() *AuthMapper {
	return &AuthMapper{}
}

func (a *AuthMapper) LoginRequestToUser(req apiModels.LoginRequest) models.User {
	return models.User{
		UserID:       req.UserID,
		VerifiedCode: req.VerifiedCode,
	}
}

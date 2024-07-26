package mappers

import (
	"github.com/wando-world/wando-sso/domain"
	apiModels "github.com/wando-world/wando-sso/internal/rest/dto"
)

type IAuthMapper interface {
	LoginRequestToUser(req apiModels.LoginRequest) domain.User
}

type AuthMapper struct {
}

func NewAuthMapper() *AuthMapper {
	return &AuthMapper{}
}

func (a *AuthMapper) LoginRequestToUser(req apiModels.LoginRequest) domain.User {
	return domain.User{
		UserID:       req.UserID,
		VerifiedCode: req.VerifiedCode,
	}
}

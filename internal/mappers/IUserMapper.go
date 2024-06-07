package mappers

import (
	apiModels "github.com/wando-world/wando-sso/internal/api/models"
	"github.com/wando-world/wando-sso/internal/models"
)

type IUserMapper interface {
	CreateUserRequestToUser(req apiModels.CreateUserRequest) (models.User, error)
}

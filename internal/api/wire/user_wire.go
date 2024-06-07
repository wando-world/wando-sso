//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/wando-world/wando-sso/internal/api"
	"github.com/wando-world/wando-sso/internal/db"
	"github.com/wando-world/wando-sso/internal/mappers"
)

func InitUserHandler() *api.UserHandler {
	wire.Build(db.NewUserRepository, mappers.NewUserMapper, api.NewUserHandler)
	return nil
}

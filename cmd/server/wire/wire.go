//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"tonotdolist/internal/app"
	"tonotdolist/internal/handler"
	"tonotdolist/internal/log"
	"tonotdolist/internal/repository"
	"tonotdolist/internal/server"
	"tonotdolist/internal/service"
	"tonotdolist/pkg/api"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRepository,
	repository.NewUserRepository,
)

var serviceSet = wire.NewSet(
	service.NewUserService,
)

var handlerSet = wire.NewSet(
	api.NewRequestResponder,
	handler.NewHandler,
	handler.NewUserHandler,
)

func InitializeApp(logger zerolog.Logger, conf *viper.Viper) *app.App {
	panic(wire.Build(
		log.NewGormLogger,
		repositorySet,
		serviceSet,
		handlerSet,

		server.NewHTTPServer,

		app.NewApp,
	))
}

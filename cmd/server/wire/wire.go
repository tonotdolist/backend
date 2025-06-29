//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"tonotdolist/internal/app"
	"tonotdolist/internal/handler"
	"tonotdolist/internal/repository"
	"tonotdolist/internal/server"
	"tonotdolist/internal/service"
	"tonotdolist/pkg/log"
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

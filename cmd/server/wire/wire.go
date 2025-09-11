//go:build wireinject
// +build wireinject

package wire

import (
	"context"
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"tonotdolist/internal/app"
	"tonotdolist/internal/handler"
	"tonotdolist/internal/repository"
	"tonotdolist/internal/server"
	"tonotdolist/internal/service"
	"tonotdolist/pkg/clock"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewRepository,
	repository.NewUserRepository,
	repository.NewSessionRepository,
	repository.NewActivityRepository,
)

var serviceSet = wire.NewSet(
	service.NewUserService,
	service.NewActivityService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewActivityHandler,
)

func InitializeApp(ctx context.Context, logger zerolog.Logger, conf *viper.Viper) *app.App {
	panic(wire.Build(
		clock.NewSystemClock,

		repositorySet,
		serviceSet,
		handlerSet,

		server.NewHTTPServer,

		app.NewApp,
	))
}

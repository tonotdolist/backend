//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"tonotdolist/internal/app"
	"tonotdolist/internal/server"
)

//var repositorySet = wire.NewSet(
//	repository.NewDB,
//	repository.NewRepository,
//	)
//
//var serviceSet = wire.NewSet()
//
//var handlerSet = wire.NewSet()

func InitializeApp(logger zerolog.Logger, conf *viper.Viper) *app.App {
	panic(wire.Build(
		//log.NewGormLogger,

		//repositorySet,
		//serviceSet,
		//handlerSet,

		server.NewHTTPServer,

		app.NewApp,
	))
}

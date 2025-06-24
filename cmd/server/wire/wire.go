//go:build wireinject
// +build wireinject

package wire

import (
	"tonotdolist/internal/app"
	"tonotdolist/internal/server"
	"tonotdolist/pkg/log"

	"github.com/google/wire"
)

var repositorySet = wire.NewSet()

var serviceSet = wire.NewSet()

var handlerSet = wire.NewSet()

func InitializeApp() *app.App {
	panic(wire.Build(
		app.NewApp,
		server.NewHTTPServer,
		log.NewLogger,
		repositorySet,
		serviceSet,
		handlerSet,
	))
}

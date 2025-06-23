package app

import (
	"context"
	"tonotdolist/pkg/server/http"
)

type App struct {
	server *http.Server
}

func NewApp(server *http.Server) *App {
	return &App{
		server: server,
	}
}

func (a *App) Start() {
	a.server.Start()
}

func (a *App) Stop(ctx context.Context) {
	err := a.server.Stop(ctx)
	if err != nil {
		return
	}
}

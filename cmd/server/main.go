package main

import (
	"context"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tonotdolist/cmd/server/wire"
	"tonotdolist/internal/log"
	"tonotdolist/pkg/config"
)

func main() {
	logger := log.NewLogger()
	conf := config.NewConfig(logger, "local.yaml")
	if conf.GetBool("prod") {
		logger = logger.Level(zerolog.WarnLevel)
	}

	app := wire.InitializeApp(logger, conf)
	app.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	app.Stop(ctx)
}

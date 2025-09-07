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
	conf, err := config.NewConfig("config/local.yaml")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load config")
	}

	if conf.GetBool("prod") {
		logger = logger.Level(zerolog.WarnLevel)
	}

	ctx := context.Background()

	app := wire.InitializeApp(ctx, logger, conf)
	app.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	stopCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	app.Stop(stopCtx)
}

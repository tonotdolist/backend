package main

import (
	"context"
	"os"
	"tonotdolist/internal/log"
	_ "tonotdolist/internal/model"
	"tonotdolist/internal/repository"
	"tonotdolist/pkg/config"
	"tonotdolist/pkg/migrate"
)

func main() {
	logger := log.NewLogger()

	f, err := os.Open("config/local.yaml")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open config")
	}
	conf, err := config.NewConfig(f)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load config")
	}

	ctx := context.Background()
	db := repository.NewDB(ctx, logger, conf)

	err = migrate.Migrate(db)

	if err != nil {
		logger.Fatal().Err(err).Msg("error migrating DB")
	}

	logger.Info().Msg("migration success")
}

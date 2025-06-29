package main

import (
	"tonotdolist/internal/log"
	_ "tonotdolist/internal/model"
	"tonotdolist/internal/repository"
	"tonotdolist/pkg/config"
	"tonotdolist/pkg/migrate"
)

func main() {
	logger := log.NewLogger()
	conf := config.NewConfig(logger, "config/local.yaml")
	db := repository.NewDB(logger, conf)

	err := migrate.Migrate(db)

	if err != nil {
		logger.Fatal().Err(err).Msg("error migrating DB")
	}

	logger.Info().Msg("migration success")
}

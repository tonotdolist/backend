package repository

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"tonotdolist/pkg/config"
)

const (
	DsnKey = "db.dsn"
)

func init() {
	config.RegisterRequiredKey(DsnKey)
}

func NewDB(logger zerolog.Logger, gormLogger logger.Interface, config *viper.Viper) *gorm.DB {
	dsn := config.GetString(DsnKey)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to connect to db")
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to get sql.DB")
	}

	if err := sqlDB.Ping(); err != nil {
		logger.Fatal().Err(err).Msg("unable to ping db")
	}

	return db
}

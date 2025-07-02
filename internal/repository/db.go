package repository

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"tonotdolist/internal/log"
	"tonotdolist/pkg/config"
)

const (
	dsnKey      = "db.dsn"
	dbDialector = "db.dialector"
)

func init() {
	config.RegisterRequiredKey(dsnKey, dbDialector)
}

func NewDB(logger zerolog.Logger, config *viper.Viper) *gorm.DB {
	dsn := config.GetString(dsnKey)
	dialectorType := config.GetString(dbDialector)

	var (
		dialector gorm.Dialector
	)

	switch strings.ToLower(dialectorType) {
	case "mysql":
		{
			dialector = mysql.Open(dsn)
		}
	case "postgres":
		{
			dialector = postgres.Open(dsn)
		}
	default:
		{
			logger.Fatal().Str("dialector", dialectorType).Msg("unable to find dialector")
		}
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: log.NewGormLogger(logger, config),
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

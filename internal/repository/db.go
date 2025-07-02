package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"time"
	"tonotdolist/internal/log"
	"tonotdolist/pkg/config"
)

const (
	dsnKey      = "db.dsn"
	dbDialector = "db.dialector"

	redisAddrKey = "cache.redis.addr"
)

func init() {
	config.RegisterRequiredKey(dsnKey, dbDialector, redisAddrKey)
}

func NewDB(ctx context.Context, logger zerolog.Logger, config *viper.Viper) *gorm.DB {
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

	withCancel, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(withCancel); err != nil {
		logger.Fatal().Err(err).Msg("unable to ping db")
	}

	return db
}

func NewRedis(ctx context.Context, logger zerolog.Logger, config *viper.Viper) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.GetString(redisAddrKey),
	})

	withCancel, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := rdb.Ping(withCancel).Err(); err != nil {
		logger.Fatal().Err(err).Msg("failed to ping redis")
		return nil
	}

	return rdb
}

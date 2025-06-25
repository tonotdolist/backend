package log

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"os"
	"tonotdolist/pkg/log"
)

func NewLogger() zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func NewGormLogger(logger zerolog.Logger, viper *viper.Viper) *log.GormLogger {
	return log.NewGormLogger(logger, viper.GetBool("db.logging.ignoreRecordNotFoundErr"), viper.GetInt64("db.logging.slowThreshold"))
}

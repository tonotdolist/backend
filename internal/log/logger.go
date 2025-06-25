package log

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"os"
	"tonotdolist/pkg/log"
)

func NewLogger(viper *viper.Viper) zerolog.Logger {
	level := zerolog.DebugLevel
	if viper.GetBool("prod") {
		level = zerolog.WarnLevel
	}

	return zerolog.New(os.Stdout).Level(level).With().Timestamp().Logger()
}

func NewGormLogger(logger zerolog.Logger, viper *viper.Viper) *log.GormLogger {
	return log.NewGormLogger(logger, viper.GetBool("db.logging.ignoreRecordNotFoundErr"), viper.GetInt64("db.logging.slowThreshold"))
}

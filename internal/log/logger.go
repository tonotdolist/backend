package log

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"os"
	"tonotdolist/pkg/config"
	"tonotdolist/pkg/log"
)

const (
	ignoreRecordNotFoundErrKey = "db.logging.ignoreRecordNotFoundErr"
	slowThresholdKey           = "db.logging.slowThreshold"
)

func init() {
	config.RegisterRequiredKey(ignoreRecordNotFoundErrKey, slowThresholdKey)
}

func NewLogger() zerolog.Logger {
	level := zerolog.DebugLevel

	return zerolog.New(os.Stdout).Level(level).With().Timestamp().Logger()
}

func NewGormLogger(logger zerolog.Logger, viper *viper.Viper) *log.GormLogger {
	return log.NewGormLogger(logger, viper.GetBool(ignoreRecordNotFoundErrKey), viper.GetInt64(slowThresholdKey))
}

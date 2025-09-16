package log

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
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
	writer := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
	return zerolog.New(writer).Level(zerolog.DebugLevel).With().Timestamp().Logger()
}

func NewGormLogger(logger zerolog.Logger, viper *viper.Viper) logger.Interface {
	return log.NewGormLogger(logger, viper.GetBool(ignoreRecordNotFoundErrKey), viper.GetInt64(slowThresholdKey))
}

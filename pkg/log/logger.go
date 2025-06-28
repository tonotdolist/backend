package log

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type GormLogger struct {
	logger zerolog.Logger

	ignoreNotFoundErr bool
	slowThreshold     int64
}

func GetLoggerFromContext(ctx context.Context) (zerolog.Logger, error) {
	rawLogger := ctx.Value("logger")
	if rawLogger == nil {
		return zerolog.Nop(), errors.New("logger not found in context")
	}
	logger1, ok := rawLogger.(zerolog.Logger)
	if !ok {
		return zerolog.Nop(), errors.New("cannot cast object in context to logger")
	}

	return logger1, nil
}

func NewGormLogger(logger zerolog.Logger, ignoreNotFoundErr bool, slowThreshold int64) *GormLogger {
	return &GormLogger{
		logger:            logger,
		ignoreNotFoundErr: ignoreNotFoundErr,
		slowThreshold:     slowThreshold,
	}
}

func (gl *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	var newLevel zerolog.Level
	switch level {
	case logger.Info:
		{
			newLevel = zerolog.InfoLevel
			break
		}
	case logger.Warn:
		{
			newLevel = zerolog.WarnLevel
			break
		}
	case logger.Error:
		{
			newLevel = zerolog.ErrorLevel
			break
		}
	case logger.Silent:
		{
			newLevel = zerolog.Disabled
			break
		}
	}

	return &GormLogger{
		logger: gl.logger.Level(newLevel),
	}
}

func (gl *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	gl.logger.Info().Ctx(ctx).Msgf(msg, args)
}

func (gl *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	gl.logger.Warn().Ctx(ctx).Msgf(msg, args)
}

func (gl *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	gl.logger.Error().Ctx(ctx).Msgf(msg, args)
}

func (gl *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin).Milliseconds()
	sql, rows := fc()

	var event *zerolog.Event

	if err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !gl.ignoreNotFoundErr) {
		event = gl.logger.Error()
	} else if elapsed > gl.slowThreshold {
		event = gl.logger.Warn()
	} else {
		event = gl.logger.Info()
	}

	event.Ctx(ctx).Err(err).Str("source", utils.FileWithLineNum()).Str("sql", sql).Int64("rows_affected", rows).Int64("elapsed_ms", elapsed).Bool("slow", elapsed > gl.slowThreshold).Msg("sql execution finished")
}

var _ logger.Interface = (*GormLogger)(nil)

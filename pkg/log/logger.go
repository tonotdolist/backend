package log

import (
	"github.com/rs/zerolog"
	"os"
)

func NewLogger() zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}

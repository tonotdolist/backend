package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RequestLogMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger = logger.With().
			Str("request_url", ctx.Request.URL.String()).
			Str("request_method", ctx.Request.Method).
			Logger()

		ctx.Set("logger", logger)

		logger.Info().Msg("request received")

		ctx.Next()
	}
}

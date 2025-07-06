package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"tonotdolist/api"
	"tonotdolist/common"
	"tonotdolist/internal/service"
	"tonotdolist/pkg/log"
)

func AuthMiddleware(userService service.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId := ctx.GetHeader("Authorization")
		if sessionId == "" {
			api.HandleResponse(ctx, common.ErrUnauthorized, nil)

			ctx.Abort()
			return
		}

		userId, err := userService.GetSession(ctx, sessionId)
		if err != nil {
			if !errors.Is(err, common.ErrUnauthorized) {
				logger := log.GetLoggerFromContext(ctx)
				logger.Error().Err(err).Msg("fail to validate session")
			}

			api.HandleResponse(ctx, err, nil)

			ctx.Abort()
			return
		}

		ctx.Set("session_id", sessionId)
		ctx.Set("user_id", userId)

		ctx.Next()
	}
}

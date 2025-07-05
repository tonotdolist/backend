package middleware

import (
	"github.com/gin-gonic/gin"
	"tonotdolist/api"
	"tonotdolist/common"
	"tonotdolist/internal/service"
)

func AuthMiddleware(userService service.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId := ctx.GetHeader("Authentication")
		if sessionId == "" {
			api.HandleResponse(ctx, common.ErrUnauthorized, nil)
			return
		}

		userId, err := userService.GetSession(ctx, sessionId)
		if err != nil {
			api.HandleResponse(ctx, err, nil)
		}

		ctx.Set("user_id", userId)
	}
}

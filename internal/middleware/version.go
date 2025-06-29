package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func VersionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("version", strings.Split(ctx.Request.URL.Path, "/"))

		ctx.Next()
	}
}

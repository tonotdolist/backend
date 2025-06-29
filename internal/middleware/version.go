package middleware

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func VersionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		versionString := strings.Replace(strings.Split(ctx.Request.URL.Path, "/")[1], "v", "", 1)
		version, _ := strconv.ParseUint(versionString, 10, 0)

		ctx.Set("version", uint(version))

		ctx.Next()
	}
}

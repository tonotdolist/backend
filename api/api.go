package api

import (
	"github.com/gin-gonic/gin"
	"tonotdolist/common"
	"tonotdolist/pkg/api"
	"tonotdolist/pkg/log"
)

func HandleResponse(ctx *gin.Context, err error, data interface{}) {
	if data == nil {
		data = map[string]string{}
	}

	apiVersion := ctx.GetUint(api.ApiVersionContextKey)

	if err == nil {
		err = common.ErrSuccess
	}

	apiHandler, ok := api.GetApiVersion(apiVersion)
	if !ok {
		ctx.JSON(500, map[string]string{})
		logger := log.GetLoggerFromContext(ctx)
		logger.Error().Msg("unknown api version")
		return
	}

	httpCode, resp := apiHandler.HandleResponse(api.GetError(apiHandler, err), data)

	ctx.JSON(httpCode, resp)
}

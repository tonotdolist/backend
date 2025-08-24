package api

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"tonotdolist/common"
	"tonotdolist/pkg/api"
	"tonotdolist/pkg/log"
)

func HandleResponse(ctx *gin.Context, err error, resp interface{}) {
	if resp == nil {
		resp = map[string]string{}
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

	versionedResp, err := api.GetResponse(resp, apiVersion)
	if err != nil {
		logger := log.GetLoggerFromContext(ctx)
		logger.Error().Msgf("unable to get versioned response: %s", err.Error())
		ctx.JSON(500, map[string]string{})
	}

	httpCode, resp := apiHandler.HandleResponse(api.GetError(apiHandler, err), versionedResp)

	ctx.JSON(httpCode, resp)
}

func BindJSON(ctx *gin.Context, internalType reflect.Type, version uint) (interface{}, error) {
	value, err := api.GetRequest(internalType, version)
	if err != nil {
		return nil, err
	}
	if err = ctx.ShouldBindJSON(value); err != nil {
		return nil, common.ErrBadRequest
	}

	return value.ToInternalRequest(), nil
}

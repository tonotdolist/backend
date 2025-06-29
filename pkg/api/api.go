package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"tonotdolist/common"
)

type RequestResponder struct {
	logger zerolog.Logger
}

type ApiVersionHandler interface {
	GetApiVersion() uint
	GetInternalError() interface{}
	HandleResponse(interface{}, interface{}) (int, interface{})
}

var versions = make(map[uint]ApiVersionHandler)

func RegisterApiVersion(handler ApiVersionHandler) {
	versions[handler.GetApiVersion()] = handler
}

func getApiVersion(version uint) (ApiVersionHandler, bool) {
	v, ok := versions[version]

	return v, ok
}

func NewRequestResponder(logger zerolog.Logger) *RequestResponder {
	return &RequestResponder{logger: logger}
}

func (rr *RequestResponder) HandleResponse(version uint, ctx *gin.Context, err error, data interface{}) {
	if data == nil {
		data = map[string]string{}
	}

	if err == nil {
		err = common.ErrSuccess
	}

	apiHandler, ok := getApiVersion(version)
	if !ok {
		ctx.JSON(500, map[string]string{})

		rr.logger.Error().Uint("api_version", version).Msg("unknown api version")
		return
	}

	httpCode, resp := apiHandler.HandleResponse(getError(apiHandler, err), data)

	ctx.JSON(httpCode, resp)
}

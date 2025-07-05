package api

import (
	"github.com/rs/zerolog"
)

const (
	ApiVersionContextKey = "version"
)

type RequestResponder struct {
	logger zerolog.Logger
}

type ApiVersionHandler interface {
	GetApiVersion() uint
	GetInternalError() interface{}
	HandleResponse(rawErr interface{}, data interface{}) (httpCode int, respData interface{})
}

var versions = make(map[uint]ApiVersionHandler)

func RegisterApiVersion(handler ApiVersionHandler) {
	versions[handler.GetApiVersion()] = handler
}

func GetApiVersion(version uint) (ApiVersionHandler, bool) {
	v, ok := versions[version]

	return v, ok
}

func NewRequestResponder(logger zerolog.Logger) *RequestResponder {
	return &RequestResponder{logger: logger}
}

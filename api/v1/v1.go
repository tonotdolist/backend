package v1

import "tonotdolist/pkg/api"

const version = 1

func init() {
	api.RegisterApiVersion(&apiVersionHandler{})
}

type apiVersionHandler struct {
}

func (a *apiVersionHandler) GetApiVersion() uint {
	return version
}

func (a *apiVersionHandler) GetInternalError() interface{} {
	return internalServerErr
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (a *apiVersionHandler) HandleResponse(rawError interface{}, data interface{}) (int, interface{}) {
	err := rawError.(*Error)
	resp := Response{Code: err.Code, Message: err.Message, Data: data}

	return err.HTTPCode, resp
}

var _ api.ApiVersionHandler = (*apiVersionHandler)(nil)

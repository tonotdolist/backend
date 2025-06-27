package v1

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HandleSuccess(ctx *gin.Context, data interface{}) {
	HandleError(ctx, ErrSuccess, data)
}

func HandleError(ctx *gin.Context, err *Error, data interface{}) {
	if data == nil {
		data = map[string]string{}
	}

	resp := Response{Code: err.Code, Message: err.Message, Data: data}
	ctx.JSON(err.HTTPCode, resp)
}

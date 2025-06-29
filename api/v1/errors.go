package v1

import (
	"tonotdolist/common"
	"tonotdolist/pkg/api"
)

func init() {
	registerError(200, 200, "Ok", common.ErrSuccess)
	registerError(400, 400, "Bad Request", common.ErrBadRequest)
	registerError(401, 40, "Unauthorized", common.ErrUnauthorized)
	registerError(409, 409, "Conflict", common.ErrConflict)
	registerError(404, 404, "Not Found", common.ErrNotFound)
}

var internalServerErr = newError(500, 500, "Internal Server Error")

type Error struct {
	HTTPCode int
	Code     int
	Message  string
}

func registerError(httpCode int, code int, msg string, internalMapping error) {
	api.RegisterError(version, newError(httpCode, code, msg), internalMapping)
}

func newError(httpCode int, code int, msg string) *Error {
	return &Error{HTTPCode: httpCode, Code: code, Message: msg}
}

package common

import "errors"

var errorList = make(map[error]struct{})

var (
	ErrSuccess      = newError("success")
	ErrBadRequest   = newError("bad request")
	ErrUnauthorized = newError("unauthorized")
	ErrNotFound     = newError("content not found")
	ErrConflict     = newError("conflict")

	ErrBadPassword      = newError("bad password")
	ErrPasswordTooShort = newError("password too short")
	ErrPasswordTooLong  = newError("password too long")
)

func newError(msg string) error {
	err := errors.New(msg)
	errorList[err] = struct{}{}
	return err
}

func IsCommonError(err error) bool {
	_, ok := errorList[err]

	return ok
}

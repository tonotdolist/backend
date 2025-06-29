package common

import "errors"

var (
	ErrSuccess      = newError("success")
	ErrBadRequest   = newError("bad request")
	ErrUnauthorized = newError("unauthorized")
	ErrNotFound     = newError("content not found")
	ErrConflict     = newError("conflict")
)

func newError(msg string) error {
	return errors.New(msg)
}

package common

import "errors"

var (
	ErrSuccess    = newError("success")
	ErrBadRequest = newError("bad request")
	ErrNotFound   = newError("content not found")
)

func newError(msg string) error {
	return errors.New(msg)
}

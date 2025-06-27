package v1

var (
	ErrSuccess             = newError(200, 200, "Ok")
	ErrBadRequest          = newError(400, 400, "Bad Request")
	ErrUnauthorized        = newError(401, 401, "Unauthorized")
	ErrNotFound            = newError(404, 404, "Not Found")
	ErrInternalServerError = newError(500, 500, "Internal Server Error")
)

type Error struct {
	HTTPCode int
	Code     int
	Message  string
}

func (e *Error) Error() string {
	return e.Message
}

func newError(httpCode int, code int, msg string) *Error {
	return &Error{httpCode, code, msg}
}

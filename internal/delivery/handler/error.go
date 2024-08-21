package handler

const (
	ErrUnknown      = "Unknown Error"
	ErrUnauthorized = "Unauthorized"
	ErrForbidden    = "Forbidden"
	ErrNotFound     = "Not Found"
	ErrConflict     = "Conflict"
	ErrInternal     = "Internal Server Error"
)

type Error interface {
	Error() string
	Code() int
	Msg() string
}

type ApiError struct {
	Err  error
	code int
	msg  string
}

func NewError(err error, code int, msg string) Error {
	return &ApiError{Err: err, code: code, msg: msg}
}

func (e *ApiError) Error() string {
	return e.Err.Error()
}

func (e *ApiError) Code() int {
	return e.code
}

func (e *ApiError) Msg() string {
	return e.msg
}

package msgerror

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
	Msg() string
}

type ApiError struct {
	Err error
	msg string
}

func NewError(err error, msg string) Error {
	return &ApiError{Err: err, msg: msg}
}

func (e *ApiError) Error() string {
	return e.Err.Error()
}

func (e *ApiError) Msg() string {
	return e.msg
}

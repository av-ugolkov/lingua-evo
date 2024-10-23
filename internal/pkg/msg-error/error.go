package msgerror

const (
	ErrUnknown      = "Unknown Error"
	ErrUnauthorized = "Unauthorized"
	ErrForbidden    = "Forbidden"
	ErrNotFound     = "Not Found"
	ErrConflict     = "Conflict"
	ErrInternal     = "Internal Server Error"
)

type ApiError struct {
	Err error
	Msg string
}

func NewError(err error, msg string) error {
	return &ApiError{Err: err, Msg: msg}
}

func (e *ApiError) Error() string {
	return e.Err.Error()
}

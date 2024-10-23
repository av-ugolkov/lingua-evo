package msgerror

const (
	ErrMsgUnknown      = "You doesn't fill one or several fields."
	ErrMsgUnauthorized = "You are not authorized."
	ErrMsgForbidden    = "Access denied."
	ErrMsgNotFound     = "Not found."
	ErrMsgConflict     = "It's already exist."
	ErrMsgInternal     = "Sorry, something went wrong."
	ErrMsgBadRequest   = "You don't fill one or several fields."
	ErrMsgBadEmail     = "Email format is invalid"
)

type ApiError struct {
	Err error
	Msg string
}

func New(err error, msg string) error {
	return &ApiError{Err: err, Msg: msg}
}

func (e *ApiError) Error() string {
	return e.Err.Error()
}

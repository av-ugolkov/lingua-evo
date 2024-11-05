package msgerr

const (
	ErrMsgUnknown      = "You doesn't fill one or several fields."
	ErrMsgUnauthorized = "You are not authorized. Update the page."
	ErrMsgForbidden    = "Access denied."
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

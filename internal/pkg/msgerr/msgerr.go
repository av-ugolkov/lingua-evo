package msgerr

const (
	ErrMsgUnknown      = "You doesn't fill one or several fields."
	ErrMsgUnauthorized = "You are not authorized. Update the page."
	ErrMsgForbidden    = "Access denied."
	ErrMsgInternal     = "Sorry, something went wrong."
	ErrMsgBadRequest   = "You don't fill one or several fields."
	ErrMsgBadEmail     = "Email format is invalid"
)

type Error struct {
	Err error
	Msg string
}

func New(err error, msg string) error {
	return &Error{Err: err, Msg: msg}
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}

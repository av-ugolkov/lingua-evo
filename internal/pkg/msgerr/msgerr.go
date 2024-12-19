package msgerr

const (
	ErrMsgUnknown      = "You doesn't fill one or several fields."
	ErrMsgUnauthorized = "You are not authorized. Update the page."
	ErrMsgForbidden    = "Access denied."
	ErrMsgInternal     = "Sorry, something went wrong."
	ErrMsgBadRequest   = "You don't fill one or several fields."
	ErrMsgBadEmail     = "Email format is invalid"
)

type MsgErr struct {
	err error
	msg string
}

func New(err error, msg string) error {
	return &MsgErr{err: err, msg: msg}
}

func (e *MsgErr) Error() string {
	return e.err.Error()
}

func (e *MsgErr) Unwrap() error {
	return e.err
}

func (e *MsgErr) Msg() string {
	return e.msg
}

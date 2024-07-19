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
	Map() map[string]any
	GetCode() int
}

type ApiError struct {
	Err  error  `json:"-"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewError(err error, code int, msg string) Error {
	return &ApiError{Err: err, Code: code, Msg: msg}
}

func (e *ApiError) Error() string {
	return e.Err.Error()
}

func (e *ApiError) Map() map[string]any {
	return map[string]any{"code": e.Code, "msg": e.Msg}
}

func (e *ApiError) GetCode() int {
	return e.Code
}

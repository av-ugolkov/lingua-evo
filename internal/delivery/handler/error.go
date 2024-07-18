package handler

import (
	"fmt"
)

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
	JSON() string
	GetCode() int
}

type ApiError struct {
	Err     error
	Code    int
	Message string
}

func NewError(err error, code int, msg string) Error {
	return &ApiError{Err: err, Code: code, Message: msg}
}

func (e *ApiError) Error() string {
	return e.Err.Error()
}

func (e *ApiError) JSON() string {
	return fmt.Sprintf(`{"code":%d,"message":"%s"}`, e.Code, e.Message)
}

func (e *ApiError) GetCode() int {
	return e.Code
}

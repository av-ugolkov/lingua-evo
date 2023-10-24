package entity

import (
	"errors"
	"fmt"
)

const (
	UsernameLen    = 3
	MinPasswordLen = 6
	PasswordSolt   = 11
)

var (
	ErrEmailNotCorrect   = errors.New("email is not correct")
	ErrItIsAdmin         = errors.New("it is admin")
	ErrEmailBusy         = errors.New("this email is busy")
	ErrUsernameLen       = fmt.Errorf("username must be more %d characters", UsernameLen)
	ErrUsernameAdmin     = errors.New("username can not contains admin")
	ErrUsernameBusy      = errors.New("this username is busy")
	ErrPasswordLen       = fmt.Errorf("password must be more %d characters", MinPasswordLen)
	ErrPasswordDifficult = errors.New("password must be more difficult")
)

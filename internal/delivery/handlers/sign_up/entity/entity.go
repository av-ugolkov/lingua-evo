package entity

import (
	"errors"
	"fmt"
)

const (
	UsernameLen = 3
	PasswordLen = 6

	UsernameAdmin = "admin"
)

var (
	ErrEmailNotCorrect   = errors.New("email is not correct")
	ErrItIsAdmin         = errors.New("it is admin")
	ErrEmailBusy         = errors.New("this email is busy")
	ErrUsernameLen       = fmt.Errorf("username must be more %d characters", UsernameLen)
	ErrUsernameAdmin     = errors.New("username can not contains admin")
	ErrUsernameBusy      = errors.New("this username is busy")
	ErrPasswordLen       = fmt.Errorf("password must be more %d characters", PasswordLen)
	ErrPasswordDifficult = errors.New("password must be more difficult")
)

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Refresh struct {
	RefreshToken string `json:"refresh_token"`
}

package entity

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	UsernameLen    = 3
	MinPasswordLen = 6
)

var (
	ErrNotFoundUser      = errors.New("not found user")
	ErrEmailNotCorrect   = errors.New("email is not correct")
	ErrItIsAdmin         = errors.New("it is admin")
	ErrEmailBusy         = errors.New("this email is busy")
	ErrUsernameLen       = fmt.Errorf("username must be more %d characters", UsernameLen)
	ErrUsernameAdmin     = errors.New("username can not contains admin")
	ErrUsernameBusy      = errors.New("this username is busy")
	ErrPasswordLen       = fmt.Errorf("password must be more %d characters", MinPasswordLen)
	ErrPasswordDifficult = errors.New("password must be more difficult")
)

type (
	User struct {
		ID           uuid.UUID
		Username     string
		Email        string
		PasswordHash string
		Role         string
		CreatedAt    time.Time
		LastVisitAt  time.Time
	}

	Session struct {
		UserID      uuid.UUID `json:"user_id"`
		Fingerprint string    `json:"fingerprint"`
	}
)

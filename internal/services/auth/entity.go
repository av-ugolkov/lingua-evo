package auth

import (
	"errors"
	"time"

	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
)

var (
	ErrEmailNotCorrect = errors.New("email is not correct")
	ErrWrongPassword   = errors.New("wrong password")
	ErrItIsAdmin       = errors.New("it is admin")
	ErrEmailBusy       = errors.New("this email is busy")
)

type (
	Session uuid.UUID

	User struct {
		ID       uuid.UUID
		Name     string
		Password string
		Email    string
		Role     runtime.Role
		Code     int
	}

	Tokens struct {
		AccessToken  string
		RefreshToken uuid.UUID
	}

	Claims struct {
		ID        uuid.UUID
		UserID    uuid.UUID
		ExpiresAt time.Time
	}
)

func (s *Session) String() string {
	u := uuid.UUID(*s)
	return u.String()
}

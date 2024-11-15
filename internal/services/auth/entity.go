package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
)

// TODO вынести в конфиги
const (
	MinNicknameLen = 3
	MinPasswordLen = 6
)

var (
	ErrNotFoundUser      = errors.New("not found user")
	ErrEmailNotCorrect   = errors.New("email is not correct")
	ErrWrongPassword     = errors.New("wrong password")
	ErrItIsAdmin         = errors.New("it is admin")
	ErrEmailBusy         = errors.New("this email is busy")
	ErrNicknameBusy      = errors.New("this nickname is busy")
	ErrNicknameLen       = fmt.Errorf("nickname must be more %d characters", MinNicknameLen)
	ErrPasswordLen       = fmt.Errorf("password must be more %d characters", MinPasswordLen)
	ErrPasswordDifficult = errors.New("password must be more difficult")
)

type (
	Session uuid.UUID

	User struct {
		ID       uuid.UUID
		Nickname string
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

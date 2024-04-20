package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/av-ugolkov/lingua-evo/runtime"
	"time"

	"github.com/google/uuid"
)

type (
	Session struct {
		UserID      uuid.UUID `json:"user_id"`
		Fingerprint string    `json:"fingerprint"`
		CreatedAt   time.Time `json:"created_at"`
	}

	User struct {
		ID       uuid.UUID
		Username string
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

// TODO вынести в конфиги
const (
	UsernameLen    = 3
	MinPasswordLen = 6
)

var (
	ErrEmailNotCorrect   = errors.New("email is not correct")
	ErrItIsAdmin         = errors.New("it is admin")
	ErrEmailBusy         = errors.New("this email is busy")
	ErrUsernameLen       = fmt.Errorf("username must be more %d characters", UsernameLen)
	ErrUsernameBusy      = errors.New("this username is busy")
	ErrPasswordLen       = fmt.Errorf("password must be more %d characters", MinPasswordLen)
	ErrPasswordDifficult = errors.New("password must be more difficult")
)

func (s *Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

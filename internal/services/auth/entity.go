package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/runtime"
	jsoniter "github.com/json-iterator/go"

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

type TypeToken string

const (
	Email  TypeToken = "email"
	Google TypeToken = "google"
)

type (
	Session struct {
		UserID       uuid.UUID
		RefreshToken string    `json:"refresh_token"`
		TypeToken    TypeToken `json:"type_token"`
		ExpiresAt    time.Time `json:"expires_at"`
	}

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
		RefreshToken string
	}
)

func (s *Session) JSON() ([]byte, error) {
	b, err := jsoniter.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("auth.Session.JSON: %w", err)
	}
	return b, nil
}

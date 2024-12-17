package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/aes"
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

	ErrMsgUserExists = "Sorry, the user exists with the same email or nickname"
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

func (s *Session) Marshal() ([]byte, error) {
	var err error
	s.RefreshToken, err = aes.EncryptAES(s.RefreshToken, config.GetConfig().AES.Key)
	if err != nil {
		return nil, fmt.Errorf("auth.Session.Marshal: %w", err)
	}

	b, err := jsoniter.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("auth.Session.Marshal: %w", err)
	}
	return b, nil
}

func (s *Session) Unmarshal(data []byte) error {
	err := jsoniter.Unmarshal(data, &s)
	if err != nil {
		return fmt.Errorf("auth.Session.Unmarshal: %w", err)
	}

	s.RefreshToken, err = aes.DecryptAES(s.RefreshToken, config.GetConfig().AES.Key)
	if err != nil {
		return fmt.Errorf("auth.Session.Unmarshal: %w", err)
	}

	return nil
}

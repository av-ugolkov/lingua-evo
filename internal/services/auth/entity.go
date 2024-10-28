package auth

import (
	"errors"
	"time"

	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

var (
	ErrEmailNotCorrect = errors.New("email is not correct")
	ErrWrongPassword   = errors.New("wrong password")
	ErrItIsAdmin       = errors.New("it is admin")
	ErrEmailBusy       = errors.New("this email is busy")
)

type (
	Session struct {
		UserID      uuid.UUID `json:"user_id"`
		Fingerprint string    `json:"fingerprint"`
		CreatedAt   time.Time `json:"created_at"`
	}

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

func (s *Session) MarshalBinary() ([]byte, error) {
	return jsoniter.Marshal(s)
}

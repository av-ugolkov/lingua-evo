package user

import (
	"errors"
	"time"

	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
)

var (
	ErrNotFoundUser = errors.New("not found user")
)

type (
	User struct {
		ID           uuid.UUID
		Name         string
		Email        string
		PasswordHash string
		Role         runtime.Role
		CreatedAt    time.Time
		LastVisitAt  time.Time
	}

	Session struct {
		UserID      uuid.UUID `json:"user_id"`
		Fingerprint string    `json:"fingerprint"`
	}
)

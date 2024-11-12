package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
)

var (
	ErrNotFoundUser = errors.New("not found user")
)

type (
	User struct {
		ID            uuid.UUID
		Nickname      string
		Email         string
		Role          runtime.Role
		MaxCountWords int
		CreatedAt     time.Time
		VisitedAt     time.Time
	}

	UserCreate struct {
		ID       uuid.UUID
		Nickname string
		Password string
		Email    string
		Role     runtime.Role
		Code     int
	}

	Session struct {
		UserID      uuid.UUID `json:"user_id"`
		Fingerprint string    `json:"fingerprint"`
	}

	UserData struct {
		UID     uuid.UUID
		Name    string
		Surname string
	}

	UserNewsletters struct {
		UID  uuid.UUID
		News bool
	}

	Subscriptions struct {
		ID             uuid.UUID
		UserID         uuid.UUID
		SubscriptionID int
		CountWord      int
		StartedAt      time.Time
		EndedAt        time.Time
	}
)

// TODO вынести в конфиги
const (
	MinUsernameLen = 3
	MinPasswordLen = 6
)

var (
	ErrEmailNotCorrect   = errors.New("email is not correct")
	ErrItIsAdmin         = errors.New("it is admin")
	ErrEmailBusy         = errors.New("this email is busy")
	ErrUsernameLen       = fmt.Errorf("username must be more %d characters", MinUsernameLen)
	ErrUsernameBusy      = errors.New("this username is busy")
	ErrPasswordLen       = fmt.Errorf("password must be more %d characters", MinPasswordLen)
	ErrPasswordDifficult = errors.New("password must be more difficult")
)

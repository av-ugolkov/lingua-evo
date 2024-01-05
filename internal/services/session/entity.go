package session

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UserID      uuid.UUID
	Fingerprint string
	CreatedAt   time.Time
}

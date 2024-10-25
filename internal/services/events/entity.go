package events

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Type      string
	CreatedAt time.Time
}

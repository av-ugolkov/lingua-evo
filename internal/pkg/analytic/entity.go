package analytics

import (
	"time"

	"github.com/google/uuid"
)

type Action struct {
	UserID    uuid.UUID `json:"uid"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
}

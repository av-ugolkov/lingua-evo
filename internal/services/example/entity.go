package example

import (
	"github.com/google/uuid"
	"time"
)

type (
	Example struct {
		ID        uuid.UUID
		Text      string
		CreatedAt time.Time
	}
)

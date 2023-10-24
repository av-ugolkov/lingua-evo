package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	Session struct {
		ID           uuid.UUID
		UserID       uuid.UUID
		RefreshToken uuid.UUID
		ExpiresAt    time.Time
		CreatedAt    time.Time
	}

	Tokens struct {
		JWT          string
		RefreshToken uuid.UUID
	}

	Claims struct {
		ID        uuid.UUID
		UserID    uuid.UUID
		ExpiresAt time.Time
	}
)

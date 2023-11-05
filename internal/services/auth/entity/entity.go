package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	Session struct {
		RefreshToken uuid.UUID
		UserID       uuid.UUID
		ExpiresAt    time.Time
		CreatedAt    time.Time
	}

	Tokens struct {
		AccessToken  string
		RefreshToken uuid.UUID
	}

	Claims struct {
		ID              uuid.UUID
		UserID          uuid.UUID
		Email           string
		HashFingerprint string
		ExpiresAt       time.Time
	}
)

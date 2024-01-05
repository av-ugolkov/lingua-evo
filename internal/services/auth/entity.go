package auth

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type (
	Session struct {
		UserID      uuid.UUID `json:"user_id"`
		Fingerprint string    `json:"fingerprint"`
		CreatedAt   time.Time `json:"created_at"`
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
	return json.Marshal(s)
}

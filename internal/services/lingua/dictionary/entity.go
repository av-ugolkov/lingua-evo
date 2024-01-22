package dictionary

import "github.com/google/uuid"

type Dictionary struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
}

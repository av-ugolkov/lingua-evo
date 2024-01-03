package dictionary

import "github.com/google/uuid"

type Dictionary struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Name   string
}

package dictionary

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrDictionaryNotFound = errors.New("dictionary not found")
)

type Dictionary struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Tags   []string  `json:"tags"`
}

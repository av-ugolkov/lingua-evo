package games

import "github.com/google/uuid"

type (
	ReviseGame struct {
		VocabID   uuid.UUID
		CountWord int
		TypeGame  string
	}
)

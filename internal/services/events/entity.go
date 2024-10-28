package events

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Payload   Payload
	CreatedAt time.Time
}

type PayloadType string

const (
	VocabCreated     PayloadType = "vocab_created"
	VocabDeleted     PayloadType = "vocab_deleted"
	VocabUpdated     PayloadType = "vocab_updated"
	VocabWordCreated PayloadType = "vocab_word_created"
	VocabWordDeleted PayloadType = "vocab_word_deleted"
	VocabWordUpdated PayloadType = "vocab_word_updated"
)

type Payload struct {
	Type PayloadType `json:"type"`
	Data any         `json:"data"`
}

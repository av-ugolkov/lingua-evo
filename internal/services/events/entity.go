package events

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PayloadType string

const (
	VocabCreated     PayloadType = "vocab_created"
	VocabDeleted     PayloadType = "vocab_deleted"
	VocabUpdated     PayloadType = "vocab_updated"
	VocabWordCreated PayloadType = "vocab_word_created"
	VocabWordDeleted PayloadType = "vocab_word_deleted"
	VocabWordUpdated PayloadType = "vocab_word_updated"
)

type (
	PayloadDataVocab struct {
		VocabID    *uuid.UUID `json:"vocab_id,omitempty"`
		VocabTitle string     `json:"vocab_title,omitempty"`
		DictWordID *uuid.UUID `json:"dict_word_id,omitempty"`
	}
)

type UserData struct {
	ID          uuid.UUID
	Name        string
	Role        string
	LastVisitAt time.Time
}

type Event struct {
	ID        uuid.UUID
	User      UserData
	Payload   Payload
	CreatedAt time.Time
	Watched   bool
}

type Payload struct {
	Type PayloadType
	Data any
}

func (p Payload) String() string {
	switch p.Type {
	case VocabCreated:
		return fmt.Sprintf("VocabCreated: %v", p.Data)
	case VocabDeleted:
		return fmt.Sprintf("VocabDeleted: %v", p.Data)
	case VocabUpdated:
		return fmt.Sprintf("VocabUpdated: %v", p.Data)
	case VocabWordCreated:
		return fmt.Sprintf("VocabWordCreated: %v", p.Data)
	case VocabWordDeleted:
		return fmt.Sprintf("VocabWordDeleted: %v", p.Data)
	case VocabWordUpdated:
		return fmt.Sprintf("VocabWordUpdated: %v", p.Data)
	default:
		return fmt.Sprintf("Unknown payload type: %v", p.Type)
	}
}

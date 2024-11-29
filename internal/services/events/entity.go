package events

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

type PayloadType string

const (
	VocabCreated     PayloadType = "vocab_created"
	VocabDeleted     PayloadType = "vocab_deleted"
	VocabUpdated     PayloadType = "vocab_updated"
	VocabRenamed     PayloadType = "vocab_renamed"
	VocabWordCreated PayloadType = "vocab_word_created"
	VocabWordDeleted PayloadType = "vocab_word_deleted"
	VocabWordUpdated PayloadType = "vocab_word_updated"
	VocabWordRenamed PayloadType = "vocab_word_renamed"
)

type (
	PayloadDataVocab struct {
		VocabID    *uuid.UUID `json:"vocab_id,omitempty"`
		VocabTitle string     `json:"vocab_title,omitempty"`
		DictWordID *uuid.UUID `json:"dict_word_id,omitempty"`
		DictWord   string     `json:"dict_word,omitempty"`
	}
)

type UserData struct {
	ID        uuid.UUID
	Nickname  string
	Role      string
	VisitedAt time.Time
}

type Event struct {
	ID        uuid.UUID
	User      UserData
	Type      PayloadType
	Payload   any
	CreatedAt time.Time
	Watched   bool
}

type EventWatched struct {
	EventID   uuid.UUID
	UserID    uuid.UUID
	WatchedAt time.Time
}

func (p Event) PayloadToMap() map[string]any {
	switch p.Type {
	case VocabWordCreated, VocabWordDeleted, VocabWordUpdated, VocabWordRenamed:
		data := p.Payload.(PayloadDataVocab)
		mp := make(map[string]any, 4)
		mp["vocab_id"] = data.VocabID
		mp["vocab_title"] = data.VocabTitle
		mp["dict_word_id"] = data.DictWordID
		mp["dict_word"] = data.DictWord
		return mp
	default:
		return map[string]any{}
	}
}

func Unmarshal(typePayload PayloadType, dataJSON []byte) (any, error) {
	switch typePayload {
	case VocabWordCreated, VocabWordDeleted, VocabWordUpdated, VocabWordRenamed:
		var data PayloadDataVocab
		err := jsoniter.Unmarshal(dataJSON, &data)
		if err != nil {
			return nil, err
		}
		return data, err
	default:
		return nil, fmt.Errorf("unknown payload type: %v", typePayload)
	}
}

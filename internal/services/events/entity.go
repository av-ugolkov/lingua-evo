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
	VocabWordCreated PayloadType = "vocab_word_created"
	VocabWordDeleted PayloadType = "vocab_word_deleted"
	VocabWordUpdated PayloadType = "vocab_word_updated"
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
	ID          uuid.UUID
	Name        string
	Role        string
	LastVisitAt time.Time
}

type Event struct {
	ID        uuid.UUID
	User      UserData
	Type      PayloadType
	Payload   any
	CreatedAt time.Time
	Watched   bool
}

func (p Event) PayloadToMap() map[string]any {
	switch p.Type {
	case VocabWordCreated, VocabWordDeleted, VocabWordUpdated:
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
	case VocabWordCreated, VocabWordDeleted, VocabWordUpdated:
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

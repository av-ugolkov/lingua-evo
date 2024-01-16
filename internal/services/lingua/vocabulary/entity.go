package vocabulary

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type (
	Word struct {
		Text          string `json:"text"`
		Pronunciation string `json:"pronunciation,omitempty"`
		LangCode      string `json:"language"`
	}

	Vocabulary struct {
		DictionaryId   uuid.UUID
		NativeWord     uuid.UUID
		TranslateWords TranslateWords
		Examples       Examples
		Tags           Tags
	}

	TranslateWords []uuid.UUID
	Examples       []uuid.UUID
	Tags           []uuid.UUID
)

func (w *TranslateWords) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &w)
}

func (w *TranslateWords) Value() (driver.Value, error) {
	return json.Marshal(w)
}

func (w *Examples) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &w)
}

func (w *Examples) Value() (driver.Value, error) {
	return json.Marshal(w)
}

func (w *Tags) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &w)
}

func (w *Tags) Value() (driver.Value, error) {
	return json.Marshal(w)
}

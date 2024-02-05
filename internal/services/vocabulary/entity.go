package vocabulary

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type (
	Word struct {
		Text          string `json:"text"`
		Pronunciation string `json:"pronunciation,omitempty"`
		LangCode      string `json:"language"`
	}

	Words []Word

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

	VocabularyWord struct {
		NativeWord     string
		TranslateWords []string
		Examples       []string
		Tags           []string
	}
)

func (w *Words) GetValues() []string {
	var values []string
	for _, word := range *w {
		values = append(values, word.Text)
	}
	return values
}

func (w *TranslateWords) Scan(value interface{}) error {
	b, ok := value.(string)
	if !ok {
		return errors.New("type assertion to string failed")
	}

	err := json.Unmarshal([]byte(b), &w)
	if err != nil {
		return fmt.Errorf("vocabulary.entity.Word.Scan: %w", err)
	}

	return nil
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

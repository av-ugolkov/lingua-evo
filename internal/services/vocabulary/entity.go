package vocabulary

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrDuplicate = errors.New("duplicate key value violates unique constraint")
)

type (
	Word struct {
		Text          string `json:"text"`
		Pronunciation string `json:"pronunciation,omitempty"`
		LangCode      string `json:"language,omitempty"`
	}

	Words []Word

	Vocabulary struct {
		DictionaryId   uuid.UUID
		NativeWord     uuid.UUID
		TranslateWords []uuid.UUID
		Examples       []uuid.UUID
		Tags           []uuid.UUID
	}

	VocabularyWord struct {
		Id             uuid.UUID
		NativeWord     Word
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

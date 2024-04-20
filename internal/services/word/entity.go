package word

import (
	"errors"
	"github.com/av-ugolkov/lingua-evo/internal/services/example"

	"github.com/av-ugolkov/lingua-evo/internal/services/word/model"

	"github.com/google/uuid"
)

var (
	ErrDuplicate = errors.New("duplicate key value violates unique constraint")
)

type (
	Words []model.VocabWord

	Word struct {
		ID             uuid.UUID
		VocabID        uuid.UUID
		NativeID       uuid.UUID
		TranslateWords []uuid.UUID
		Examples       []uuid.UUID
	}

	DataWord struct {
		ID            uuid.UUID
		Text          string
		Pronunciation string
	}

	VocabWord struct {
		VocabID        uuid.UUID
		WordID         uuid.UUID
		NativeWord     DataWord
		TranslateWords []DataWord
		Examples       []example.Example
	}

	VocabularyWord struct {
		ID             uuid.UUID
		NativeWord     model.VocabWord
		TranslateWords []string
		Examples       []string
	}
)

func (w *Words) GetValues() []string {
	var values []string
	for _, word := range *w {
		values = append(values, word.Text)
	}
	return values
}

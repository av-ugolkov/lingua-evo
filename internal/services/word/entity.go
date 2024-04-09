package word

import (
	"errors"

	"github.com/av-ugolkov/lingua-evo/internal/services/word/model"

	"github.com/google/uuid"
)

var (
	ErrDuplicate = errors.New("duplicate key value violates unique constraint")
)

type (
	Words []model.Word

	Word struct {
		ID             uuid.UUID
		VocabID        uuid.UUID
		NativeID       uuid.UUID
		TranslateWords []uuid.UUID
		Examples       []uuid.UUID
	}

	VocabularyWord struct {
		Id             uuid.UUID
		NativeWord     model.Word
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

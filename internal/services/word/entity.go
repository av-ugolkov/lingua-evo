package word

import (
	"errors"
	"time"

	entityDict "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	entityExample "github.com/av-ugolkov/lingua-evo/internal/services/example"

	"github.com/google/uuid"
)

var (
	ErrDuplicate         = errors.New("duplicate key value violates unique constraint")
	ErrWordPronunciation = errors.New("word pronunciation is empty")
)

type (
	VocabWord struct {
		ID            uuid.UUID
		VocabID       uuid.UUID
		NativeID      uuid.UUID
		Pronunciation string
		TranslateIDs  []uuid.UUID
		ExampleIDs    []uuid.UUID
		UpdatedAt     time.Time
		CreatedAt     time.Time
	}

	VocabWordData struct {
		ID         uuid.UUID
		VocabID    uuid.UUID
		Native     entityDict.DictWord
		Translates []entityDict.DictWord
		Examples   []entityExample.Example
		UpdatedAt  time.Time
		CreatedAt  time.Time
	}
)

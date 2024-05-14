package word

import (
	"errors"
	"time"

	entityDict "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	entityExample "github.com/av-ugolkov/lingua-evo/internal/services/example"

	"github.com/google/uuid"
)

var (
	ErrDuplicate = errors.New("duplicate key value violates unique constraint")
)

type (
	VocabWord struct {
		ID           uuid.UUID
		VocabID      uuid.UUID
		NativeID     uuid.UUID
		TranslateIDs []uuid.UUID
		ExampleIDs   []uuid.UUID
		UpdatedAt    time.Time
		CreatedAt    time.Time
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

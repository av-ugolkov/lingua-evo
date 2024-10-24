package vocabulary

import (
	"errors"
	"time"

	entityDict "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	entityExample "github.com/av-ugolkov/lingua-evo/internal/services/example"
	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	"github.com/av-ugolkov/lingua-evo/runtime/access"

	"github.com/google/uuid"
)

var (
	ErrVocabularyNotFound = errors.New("vocabulary not found")
	ErrAccessDenied       = errors.New("access denied")

	ErrDuplicate         = errors.New("duplicate key value violates unique constraint")
	ErrWordPronunciation = errors.New("Pronunciation not found")
	ErrUserWordLimit     = errors.New("user word limit reached")
)

type (
	Vocab struct {
		ID            uuid.UUID
		UserID        uuid.UUID
		Name          string
		Access        uint8
		NativeLang    string
		TranslateLang string
		Description   string
		Tags          []entityTag.Tag
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	VocabWithUser struct {
		Vocab
		UserName     string
		Editable     bool
		Notification bool
		WordsCount   uint
	}

	VocabWithUserAndWords struct {
		VocabWithUser
		Words []string
	}

	VocabWord struct {
		ID            uuid.UUID
		VocabID       uuid.UUID
		NativeID      uuid.UUID
		Pronunciation string
		Description   string
		TranslateIDs  []uuid.UUID
		ExampleIDs    []uuid.UUID
		UpdatedAt     time.Time
		CreatedAt     time.Time
	}

	VocabWordData struct {
		ID          uuid.UUID
		VocabID     uuid.UUID
		Native      entityDict.DictWord
		Description string
		Translates  []entityDict.DictWord
		Examples    []entityExample.Example
		UpdatedAt   time.Time
		CreatedAt   time.Time
	}

	Access struct {
		ID      int
		VocabID uuid.UUID
		UserID  uuid.UUID
		Status  access.Status
	}
)

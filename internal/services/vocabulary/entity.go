package vocabulary

import (
	"errors"
	"time"

	"github.com/av-ugolkov/lingua-evo/runtime/access"

	"github.com/google/uuid"
)

var (
	ErrVocabularyNotFound = errors.New("vocabulary not found")
	ErrAccessDenied       = errors.New("access denied")
	ErrDuplicate          = errors.New("duplicate key value violates unique constraint")
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
		Definition    string
		TranslateIDs  []uuid.UUID
		ExampleIDs    []uuid.UUID
		UpdatedAt     time.Time
		CreatedAt     time.Time
	}

	DictWord struct {
		ID            uuid.UUID
		Text          string
		Pronunciation string
		LangCode      string
		Creator       uuid.UUID
	}

	Example struct {
		ID   uuid.UUID
		Text string
	}

	VocabWordData struct {
		ID         uuid.UUID
		VocabID    uuid.UUID
		Native     DictWord
		Definition string
		Translates []DictWord
		Examples   []Example
		UpdatedAt  time.Time
		CreatedAt  time.Time
	}

	Access struct {
		ID      int
		VocabID uuid.UUID
		UserID  uuid.UUID
		Status  access.Status
	}
)

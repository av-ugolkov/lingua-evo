package vocabulary

import (
	"errors"
	"time"

	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	"github.com/av-ugolkov/lingua-evo/runtime/access"

	"github.com/google/uuid"
)

var (
	ErrVocabularyNotFound = errors.New("vocabulary not found")
	ErrAccessDenied       = errors.New("access denied")
)

type (
	Vocabulary struct {
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

	VocabularyWithUser struct {
		Vocabulary
		UserName   string
		WordsCount uint
	}

	Access struct {
		ID      int
		VocabID uuid.UUID
		UserID  uuid.UUID
		Status  access.Status
	}
)

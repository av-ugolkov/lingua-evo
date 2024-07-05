package vocabulary

import (
	"errors"

	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"

	"github.com/google/uuid"
)

var (
	ErrVocabularyNotFound = errors.New("vocabulary not found")
)

type (
	Vocabulary struct {
		ID            uuid.UUID
		UserID        uuid.UUID
		Name          string
		Access        int
		NativeLang    string
		TranslateLang string
		Tags          []entityTag.Tag
	}
)

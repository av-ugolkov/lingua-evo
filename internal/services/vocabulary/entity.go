package vocabulary

import (
	"errors"
	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"

	"github.com/google/uuid"
)

var (
	ErrVocabularyNotFound = errors.New("vocabulary not found")
	ErrCountVocabulary    = errors.New("too much dictionaries for user")
)

type (
	Vocabulary struct {
		ID            uuid.UUID
		UserID        uuid.UUID
		Name          string
		NativeLang    string
		TranslateLang string
		Tags          []entityTag.Tag
	}
)

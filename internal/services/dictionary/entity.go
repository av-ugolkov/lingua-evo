package dictionary

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAffectRows     = errors.New("more than 1 affected rows have changed")
	ErrDuplicateWords = errors.New("duplicate words")

	ErrMsgLanguageNotFound          = "Sorry, language not found."
	ErrMsgWordPronunciationNotFound = "Sorry, pronunciation not found"
)

type DictWord struct {
	ID            uuid.UUID
	Text          string
	Pronunciation string
	LangCode      string
	Creator       uuid.UUID
	Moderator     uuid.UUID
	UpdatedAt     time.Time
	CreatedAt     time.Time
}

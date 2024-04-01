package dictionary

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrDictionaryNotFound = errors.New("dictionary not found")
)

type Dictionary struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Name       string
	NativeLang string
	SecondLang string
	Tags       []string
}

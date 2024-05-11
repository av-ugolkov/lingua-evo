package dictionary

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrorAffectRows = errors.New("more than 1 affected rows have changed")
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

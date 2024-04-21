package dictionary

import (
	"time"

	"github.com/google/uuid"
)

type DictWord struct {
	ID            uuid.UUID
	Text          string
	Pronunciation string
	LangCode      string
	Creator       uuid.UUID
	Moderator     uuid.UUID
	UpdateAt      time.Time
	CreatedAt     time.Time
}

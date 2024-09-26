package subscribers

import (
	"time"

	"github.com/google/uuid"
)

type (
	VocabularyNotification struct {
		UserID    uuid.UUID
		VocabID   uuid.UUID
		CreatedAt time.Time
	}

	Notification struct {
		ID        uuid.UUID
		UserID    uuid.UUID
		Title     string
		Message   string
		IsRead    bool
		CreatedAt time.Time
	}
)

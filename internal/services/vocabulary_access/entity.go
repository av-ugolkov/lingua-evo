package vocabulary_access

import "github.com/google/uuid"

type AccessType int

const (
	Private AccessType = iota
	SomeOne
	Subscritebers
	Public
)

type (
	Access struct {
		ID         int
		VocabID    uuid.UUID
		UserID     uuid.UUID
		AccessEdit bool
	}
)

package vocabulary_access

import "github.com/google/uuid"

const (
	Private = iota
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

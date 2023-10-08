package entity

import "github.com/google/uuid"

type Vocabulary struct {
	DictionaryId  uuid.UUID
	OriginalWord  uuid.UUID
	TranslateWord []uuid.UUID
	Example       []uuid.UUID
	Tags          []uuid.UUID
}

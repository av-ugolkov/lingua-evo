package entity

import "github.com/google/uuid"

type Vocabulary struct {
	DictionaryId  uuid.UUID
	NativeWord    uuid.UUID
	TranslateWord []uuid.UUID
	Examples      []uuid.UUID
	Tags          []uuid.UUID
}

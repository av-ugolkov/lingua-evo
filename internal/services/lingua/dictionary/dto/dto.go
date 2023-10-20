package dto

import "github.com/google/uuid"

type (
	DictionaryRq struct {
		UserID uuid.UUID `json:"user_id"`
		Name   string    `json:"name"`
	}

	DictionariesRs struct {
		Dictionaries []Dictionary
	}

	Dictionary struct {
		ID   uuid.UUID
		Name string
	}

	DictionaryIDRs struct {
		ID uuid.UUID `json:"dictionary_id"`
	}
)

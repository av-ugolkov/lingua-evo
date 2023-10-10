package dto

import "github.com/google/uuid"

type (
	DictionaryRequest struct {
		UserID uuid.UUID
		Name   string
	}

	DictionariesResponse struct {
		Dictionaries []Dictionary
	}

	Dictionary struct {
		ID   uuid.UUID
		Name string
	}
)

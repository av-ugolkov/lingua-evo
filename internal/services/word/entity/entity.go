package entity

import "github.com/google/uuid"

type Word struct {
	ID            uuid.UUID
	Text          string
	Pronunciation string
	LanguageCode  string
}

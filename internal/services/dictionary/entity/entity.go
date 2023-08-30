package entity

import "github.com/google/uuid"

type Dictionary struct {
	UserId        uuid.UUID
	OriginalWord  uuid.UUID
	OriginalLang  string
	TranslateLang string
	TranslateWord []uuid.UUID
	Pronunciation string
	Example       []uuid.UUID
}

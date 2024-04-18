package model

import "github.com/google/uuid"

type (
	VocabularyRq struct {
		Name          string   `json:"name"`
		NativeLang    string   `json:"native_lang"`
		TranslateLang string   `json:"translate_lang"`
		Tags          []string `json:"tags"`
	}

	VocabularyIDRs struct {
		ID uuid.UUID `json:"id"`
	}

	VocabularyRs struct {
		ID            uuid.UUID `json:"id"`
		UserID        uuid.UUID `json:"user_id"`
		Name          string    `json:"name"`
		NativeLang    string    `json:"native_lang"`
		TranslateLang string    `json:"translate_lang"`
		Tags          []string  `json:"tags"`
	}
)

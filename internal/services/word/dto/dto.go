package dto

import "github.com/google/uuid"

type (
	AddWordRequest struct {
		Text          string `json:"text"`
		LanguageCode  string `json:"language_code"`
		Pronunciation string `json:"pronunciation,omitempty"`
	}

	GetWordRequest struct {
		Text         string `json:"text"`
		LanguageCode string `json:"language_code"`
	}
	GetWordIDRequest struct {
		ID uuid.UUID `json:"id"`
	}

	GetRandomWordRequest struct {
		LanguageCode string `json:"language_code"`
	}
)

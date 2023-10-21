package dto

import "github.com/google/uuid"

type (
	AddWordRq struct {
		Text          string `json:"text"`
		LanguageCode  string `json:"language_code"`
		Pronunciation string `json:"pronunciation,omitempty"`
	}

	GetWordRq struct {
		Text         string `json:"text"`
		LanguageCode string `json:"language_code"`
	}
	GetWordIDRq struct {
		ID uuid.UUID `json:"id"`
	}

	RandomWordRq struct {
		LanguageCode string `json:"language_code"`
	}

	RandomWordRs struct {
		Text          string `json:"text"`
		LanguageCode  string `json:"language_code"`
		Pronunciation string `json:"pronunciation,omitempty"`
	}
)

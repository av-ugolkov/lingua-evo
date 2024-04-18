package model

import "github.com/google/uuid"

type (
	WordRq struct {
		Text          string `json:"text"`
		LangCode      string `json:"lang_code"`
		Pronunciation string `json:"pronunciation,omitempty"`
	}

	WordIDRq struct {
		ID uuid.UUID `json:"id"`
	}

	WordRs struct {
		Text          string `json:"text"`
		LangCode      string `json:"lang_code"`
		Pronunciation string `json:"pronunciation,omitempty"`
	}

	WordIDRs struct {
		ID uuid.UUID `json:"id"`
	}
)

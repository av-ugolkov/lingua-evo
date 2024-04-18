package model

import (
	"github.com/google/uuid"
)

type (
	VocabWord struct {
		Text          string `json:"text"`
		Pronunciation string `json:"pronunciation,omitempty"`
	}

	VocabWordRq struct {
		VocabID       uuid.UUID `json:"vocab_id"`
		WordID        uuid.UUID `json:"word_id,omitempty"`
		NativeWord    VocabWord `json:"native_word"`
		TanslateWords []string  `json:"translate_words,omitempty"`
		Examples      []string  `json:"examples,omitempty"`
	}

	RemoveVocabWordRq struct {
		VocabID uuid.UUID `json:"vocab_id"`
		WordID  uuid.UUID `json:"word_id"`
	}

	VocabWordsRs struct {
		WordID         uuid.UUID  `json:"word_id"`
		NativeWord     *VocabWord `json:"native,omitempty"`
		TranslateWords []string   `json:"translate_words,omitempty"`
		Examples       []string   `json:"examples,omitempty"`
	}
)

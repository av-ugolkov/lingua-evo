package model

import (
	"github.com/google/uuid"
)

type (
	Word struct {
		Text          string `json:"text"`
		Pronunciation string `json:"pronunciation,omitempty"`
		LangCode      string `json:"language,omitempty"`
	}

	VocabWordRq struct {
		VocabID       uuid.UUID `json:"vocab_id"`
		WordID        uuid.UUID `json:"word_id,omitempty"`
		NativeWord    Word      `json:"native_word"`
		TanslateWords []string  `json:"translate_words,omitempty"`
		Examples      []string  `json:"examples,omitempty"`
	}

	RemoveVocabWordRq struct {
		VocabID uuid.UUID `json:"vocab_id"`
		WordID  uuid.UUID `json:"word_id"`
	}

	VocabWordsRs struct {
		WordID         uuid.UUID `json:"word_id,omitempty"`
		NativeWord     Word      `json:"native"`
		TranslateWords []string  `json:"translate_words,omitempty"`
		Examples       []string  `json:"examples,omitempty"`
	}
)

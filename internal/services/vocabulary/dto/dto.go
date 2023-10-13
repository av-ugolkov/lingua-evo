package dto

import "github.com/google/uuid"

type (
	AddWordRq struct {
		DictionaryID  uuid.UUID        `json:"dictionary_id"`
		NativeWord    VocabularyWord   `json:"native_word"`
		TanslateWords []VocabularyWord `json:"translate_words"`
		Examples      []string         `json:"examples"`
		Tags          []string         `json:"tags"`
	}

	VocabularyWord struct {
		Text          string `json:"text"`
		Pronunciation string `json:"pronunciation,omitempty"`
		LangCode      string `json:"language"`
	}
)

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

	RemoveWordRq struct {
		DictionaryID uuid.UUID `json:"dictionary_id"`
		NativeWordID uuid.UUID `json:"native_word_id"`
	}

	VocabularyWord struct {
		Text          string `json:"text"`
		Pronunciation string `json:"pronunciation,omitempty"`
		LangCode      string `json:"language"`
	}
)

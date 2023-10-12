package dto

import "github.com/google/uuid"

type (
	AddWordRq struct {
		DictionaryID uuid.UUID      `json:"dictionary_id"`
		OriginalWord VocabularyWord `json:"orig_word"`
		TanslateWord VocabularyWord `json:"translate_word"`
		Examples     []string       `json:"examples"`
		Tags         []string       `json:"tags"`
	}

	VocabularyWord struct {
		Text          string `json:"text"`
		Pronunciation string `json:"pronunciation,omitempty"`
		LangCode      string `json:"language"`
	}
)

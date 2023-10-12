package dto

type (
	AddWordRq struct {
		DictionaryID string         `json:"dictionary_id"`
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

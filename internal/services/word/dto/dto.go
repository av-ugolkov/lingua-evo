package dto

type (
	AddWordRequest struct {
		Text          string `json:"text"`
		Language      string `json:"language"`
		Pronunciation string `json:"pronunciation,omitempty"`
	}

	GetWordRequest struct {
		Text     string `json:"text"`
		Language string `json:"language"`
	}

	GetRandomWordRequest struct {
		Language string `json:"language"`
	}
)

package dto

type GetRandomWordRequest struct {
	Language string `json:"language"`
}

type AddWordRequest struct {
	Text          string `json:"text"`
	Language      string `json:"language"`
	Pronunciation string `json:"pronunciation,omitempty"`
}

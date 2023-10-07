package dto

type GetRandomWordRequest struct {
	Language string `json:"language"`
}

type AddWordRequest struct {
	OrigWord      string `json:"orig_word"`
	OrigLang      string `json:"orig_lang"`
	TranWord      string `json:"tran_word"`
	TranLang      string `json:"tran_lang"`
	Example       string `json:"example"`
	Pronunciation string `json:"pronunciation"`
}

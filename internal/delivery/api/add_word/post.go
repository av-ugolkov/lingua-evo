package add_word

import (
	"encoding/json"
	"errors"
	"net/http"

	"lingua-evo/internal/delivery/repository"
)

type AddWord struct {
	OrigWord string `json:"orig_word"`
	OrigLang string `json:"orig_lang"`
	TranWord string `json:"tran_word"`
	TranLang string `json:"tran_lang"`
}

func (h *Handler) postAddWord(w http.ResponseWriter, r *http.Request) {
	var data AddWord

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.logger.Error(err)
		return
	}
	defer r.Body.Close()

	if data.OrigWord == "" || data.TranWord == "" {
		h.logger.Error(errors.New("empty word can't add in db"))
		return
	}

	_, err := h.lingua.AddWord(r.Context(), &repository.Word{
		Text:     data.OrigWord,
		Language: data.OrigLang,
	})
	if err != nil {
		h.logger.Error(err)
		return
	}

	_, err = h.lingua.AddWord(r.Context(), &repository.Word{
		Text:     data.TranWord,
		Language: data.TranLang,
	})
	if err != nil {
		h.logger.Error(err)
		return
	}

	//err = h.lingua.AddWordInDictionary(ctx, "", wordId, wordId)
	//if err != nil {
	//	return uuid.UUID{}, err
	//}
}

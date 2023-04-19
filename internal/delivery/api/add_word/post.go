package add_word

import (
	"encoding/json"
	"errors"
	"fmt"
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if data.OrigWord == "" || data.TranWord == "" {
		h.logger.Error(errors.New("empty word can't add in db"))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("can't add empty"))
		return
	}

	_, err := h.lingua.SendWord(r.Context(), &repository.Word{
		Text:     data.OrigWord,
		Language: data.OrigLang,
	})
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.lingua.SendWord(r.Context(), &repository.Word{
		Text:     data.TranWord,
		Language: data.TranLang,
	})
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//err = h.lingua.AddWordInDictionary(ctx, "", wordId, wordId)
	//if err != nil {
	//	return uuid.UUID{}, err
	//}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%s => %s", data.OrigWord, data.TranWord)))
}

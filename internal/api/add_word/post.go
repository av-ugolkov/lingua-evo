package add_word

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"lingua-evo/internal/delivery/repository"
)

type AddWord struct {
	Text string `json:"text"`
	Lang string `json:"lang"`
}

func (h *Handler) postAddWord(w http.ResponseWriter, r *http.Request) {
	var data AddWord

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if data.Text == "" || data.Lang == "" {
		h.logger.Error(errors.New("empty word can't add in db"))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("can't add empty word"))
		return
	}

	word, err := h.lingua.SendWord(r.Context(), &repository.Word{
		Text:     data.Text,
		Language: data.Lang,
	})
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%v %s %s", word, data.Lang, data.Text)))
}

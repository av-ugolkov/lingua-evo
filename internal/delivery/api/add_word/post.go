package add_word

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"

	"lingua-evo/internal/delivery/repository"
)

type AddWord struct {
	OrigWord      string `json:"orig_word"`
	OrigLang      string `json:"orig_lang"`
	TranWord      string `json:"tran_word"`
	TranLang      string `json:"tran_lang"`
	Example       string `json:"example"`
	Pronunciation string `json:"pronunciation"`
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

	origWordId, err := h.lingua.AddWord(r.Context(), &repository.Word{
		Text:     data.OrigWord,
		Language: data.OrigLang,
	})
	if err != nil {
		h.logger.Error(err)
		return
	}

	var exampleId uuid.UUID
	if len(data.Example) > 0 {
		exampleId, err = h.lingua.AddExample(r.Context(), origWordId, data.Example)
		if err != nil {
			h.logger.Error(err)
			return
		}
	}

	tranWordId, err := h.lingua.AddWord(r.Context(), &repository.Word{
		Text:     data.TranWord,
		Language: data.TranLang,
	})
	if err != nil {
		h.logger.Error(err)
		return
	}

	id, err := h.lingua.AddWordInDictionary(r.Context(), uuid.Nil, origWordId, []uuid.UUID{tranWordId}, data.Pronunciation, []uuid.UUID{exampleId})
	if err != nil {
		h.logger.Error(err)
		return
	}
	h.logger.Println(id)
}

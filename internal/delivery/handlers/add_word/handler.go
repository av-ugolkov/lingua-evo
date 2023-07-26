package add_word

import (
	"encoding/json"
	"errors"
	"net/http"

	"lingua-evo/internal/delivery/handlers/add_word/entity"
	"lingua-evo/internal/delivery/repository"
	"lingua-evo/internal/service"
	staticFiles "lingua-evo/static"

	"lingua-evo/pkg/logging"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

const (
	addWordURL = "/add_word"

	addWordPage = "web/dictionary/add_word/add_word.html"
)

type Handler struct {
	logger *logging.Logger
	lingua *service.Lingua
}

func Create(log *logging.Logger, ling *service.Lingua, r *httprouter.Router) {
	handler := newHandler(log, ling)
	handler.register(r)
}

func newHandler(logger *logging.Logger, lingua *service.Lingua) *Handler {
	return &Handler{
		logger: logger,
		lingua: lingua,
	}
}

func (h *Handler) register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, addWordURL, h.get)
	router.HandlerFunc(http.MethodPost, addWordURL, h.post)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	t, err := staticFiles.ParseFiles(addWordPage)
	if err != nil {
		h.logger.Errorf("add_word.get.OpenFile: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	languages, err := h.lingua.GetLanguages(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	data := struct {
		Languages []*service.Language
	}{
		Languages: languages,
	}

	err = t.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	var data entity.AddWord

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

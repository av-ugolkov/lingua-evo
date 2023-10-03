package add_word

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	entityLanguage "lingua-evo/internal/services/language/entity"
	"lingua-evo/internal/services/word/entity"
	staticFiles "lingua-evo/static"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	addWord       = "/add_word"
	getWord       = "/get_word/{key}"
	getRandomWord = "/get_random_word"

	addWordPage = "web/dictionary/add_word/add_word.html"
)

type (
	langSvc interface {
		GetLanguages(context.Context) ([]*entityLanguage.Language, error)
	}

	wordSvc interface {
		AddWord(ctx context.Context, word *entity.Word) (uuid.UUID, error)
	}

	Handler struct {
		wordSvc wordSvc
		langSvc langSvc
	}
)

func Create(r *mux.Router, wordSvc wordSvc, langSvc langSvc) {
	handler := newHandler(wordSvc, langSvc)
	handler.register(r)
}

func newHandler(wordSvc wordSvc, langSvc langSvc) *Handler {
	return &Handler{
		wordSvc: wordSvc,
		langSvc: langSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(getWord, h.getWord).Methods(http.MethodGet)
	r.HandleFunc(getRandomWord, h.getRandomWord).Methods(http.MethodGet)
	r.HandleFunc(addWord, h.addWord).Methods(http.MethodPost)
}

/*func (h *Handler) openPage(w http.ResponseWriter, r *http.Request) {
	t, err := staticFiles.ParseFiles(addWordPage)
	if err != nil {
		slog.Error("add_word.get.OpenFile: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	languages, err := h.langSvc.GetLanguages(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	data := struct {
		Languages []*entityLanguage.Language
	}{
		Languages: languages,
	}

	err = t.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
}*/

func (h *Handler) getWord(w http.ResponseWriter, r *http.Request) {
	slog.Info(r.URL.Path)
}

func (h *Handler) getRandomWord(w http.ResponseWriter, r *http.Request) {
	t, err := staticFiles.ParseFiles(addWordPage)
	if err != nil {
		slog.Error("add_word.get.OpenFile: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	languages, err := h.langSvc.GetLanguages(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	data := struct {
		Languages []*entityLanguage.Language
	}{
		Languages: languages,
	}

	err = t.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
}

func (h *Handler) addWord(w http.ResponseWriter, r *http.Request) {
	var data entity.AddWord

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		slog.Error(err.Error())
		return
	}
	defer r.Body.Close()

	if data.OrigWord == "" || data.TranWord == "" {
		slog.Error("empty word can't add in db")
		return
	}

	origWordId, err := h.wordSvc.AddWord(r.Context(), &entity.Word{
		Text:     data.OrigWord,
		Language: data.OrigLang,
	})
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info(origWordId.String())
}

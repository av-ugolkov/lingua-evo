package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	entityLanguage "lingua-evo/internal/services/language/entity"
	"lingua-evo/internal/services/word/entity"
	"lingua-evo/pkg/tools"
	staticFiles "lingua-evo/static"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	openPage      = "/word"
	addWord       = "/add_word"
	getRandomWord = "/get_random_word"

	addWordPage = "web/dictionary/add_word/add_word.html"
)

type (
	langSvc interface {
		GetLanguages(context.Context) ([]*entityLanguage.Language, error)
		CheckLanguage(ctx context.Context, lang string) error
	}

	wordSvc interface {
		AddWord(ctx context.Context, word *entity.Word) (uuid.UUID, error)
		GetRandomWord(ctx context.Context, lang string) (*entity.Word, error)
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
	r.HandleFunc(openPage, h.openPage).Methods(http.MethodGet)
	r.HandleFunc(getRandomWord, h.getRandomWord).Methods(http.MethodPost)
	r.HandleFunc(addWord, h.addWord).Methods(http.MethodPost)
}

func (h *Handler) openPage(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) getRandomWord(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(errors.New("body is empty").Error()))
		return
	}

	ctx := r.Context()
	var randomWord GetRandomWordRequest
	defer func() {
		_ = r.Body.Close()
	}()

	if err := tools.CheckBody(w, r, &randomWord); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	if err := h.langSvc.CheckLanguage(ctx, randomWord.Language); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	word, err := h.wordSvc.GetRandomWord(r.Context(), randomWord.Language)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	_, _ = w.Write([]byte(word.Text))
}

func (h *Handler) addWord(w http.ResponseWriter, r *http.Request) {
	var data AddWordRequest

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

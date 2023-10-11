package delivery

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	entityLanguage "lingua-evo/internal/services/language/entity"
	"lingua-evo/internal/services/word/dto"
	"lingua-evo/internal/services/word/entity"
	"lingua-evo/internal/tools"
	staticFiles "lingua-evo/static"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	openPage      = "/word/openPage"
	addWord       = "/word/add"
	getWord       = "/word/get"
	getRandomWord = "/word/get_random"
)

const addWordPage = "web/dictionary/add_word/add_word.html"

type (
	langSvc interface {
		GetAvailableLanguages(ctx context.Context) ([]*entityLanguage.Language, error)
		GetLanguage(ctx context.Context, lang string) (*entityLanguage.Language, error)
		CheckLanguage(ctx context.Context, lang string) error
	}

	wordSvc interface {
		AddWord(ctx context.Context, word *entity.Word) (uuid.UUID, error)
		GetWord(ctx context.Context, text, language string) (uuid.UUID, error)
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
	r.HandleFunc(addWord, h.addWord).Methods(http.MethodPost)
	r.HandleFunc(getWord, h.getWord).Methods(http.MethodPost)
	r.HandleFunc(getRandomWord, h.getRandomWord).Methods(http.MethodPost)
}

func (h *Handler) openPage(w http.ResponseWriter, r *http.Request) {
	t, err := staticFiles.ParseFiles(addWordPage)
	if err != nil {
		slog.Error(fmt.Errorf("add_word.get.OpenFile: %v", err).Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	languages, err := h.langSvc.GetAvailableLanguages(r.Context())
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
	var data dto.AddWordRequest
	defer func() {
		_ = r.Body.Close()
	}()

	if err := tools.CheckBody(w, r, &data); err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.addWord - check body: %v", err))
		return
	}

	ctx := r.Context()

	lang, err := h.langSvc.GetLanguage(ctx, data.Language)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.addWord - get language: %v", err))
		return
	}

	word := &entity.Word{
		ID:            uuid.New(),
		Text:          data.Text,
		Pronunciation: data.Pronunciation,
		LanguageCode:  lang.Code,
	}

	wordUUID, err := h.wordSvc.AddWord(ctx, word)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, err)
		return
	}
	_, _ = w.Write([]byte(wordUUID.String()))
}

func (h *Handler) getWord(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	var data dto.GetWordRequest

	if err := tools.CheckBody(w, r, &data); err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWord - check body: %v", err))
		return
	}

	ctx := r.Context()

	if err := h.langSvc.CheckLanguage(ctx, data.Language); err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWord - check language: %v", err))
		return
	}

	wordID, err := h.wordSvc.GetWord(ctx, data.Text, data.Language)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWord: %v", err))
		return
	}
	_, _ = w.Write([]byte(wordID.String()))
}

func (h *Handler) getRandomWord(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	var data dto.GetRandomWordRequest

	if err := tools.CheckBody(w, r, &data); err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getRandomWord - check body: %v", err))
		return
	}

	ctx := r.Context()

	if err := h.langSvc.CheckLanguage(ctx, data.Language); err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getRandomWord - check language: %v", err))
		return
	}

	word, err := h.wordSvc.GetRandomWord(ctx, data.Language)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getRandomWord: %v", err))
		return
	}
	_, _ = w.Write([]byte(word.Text))
}

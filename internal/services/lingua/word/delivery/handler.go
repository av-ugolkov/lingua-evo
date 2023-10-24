package delivery

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	entityLanguage "lingua-evo/internal/services/lingua/language/entity"
	serviceLang "lingua-evo/internal/services/lingua/language/service"
	"lingua-evo/internal/services/lingua/word/dto"
	serviceWord "lingua-evo/internal/services/lingua/word/service"

	staticFiles "lingua-evo"
	"lingua-evo/pkg/tools"

	"github.com/gorilla/mux"
)

const (
	openPage      = "/word/openPage"
	addWord       = "/word/add"
	getWord       = "/word/get"
	getRandomWord = "/word/get_random"
)

const addWordPage = "dictionary/add_word/add_word.html"

type (
	Handler struct {
		wordSvc *serviceWord.WordSvc
		langSvc *serviceLang.LanguageSvc
	}
)

func Create(r *mux.Router, wordSvc *serviceWord.WordSvc, langSvc *serviceLang.LanguageSvc) {
	handler := newHandler(wordSvc, langSvc)
	handler.register(r)
}

func newHandler(wordSvc *serviceWord.WordSvc, langSvc *serviceLang.LanguageSvc) *Handler {
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
	var data dto.AddWordRq
	defer func() {
		_ = r.Body.Close()
	}()

	if err := tools.CheckBody(w, r, &data); err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.addWord - check body: %v", err))
		return
	}

	ctx := r.Context()

	err := h.langSvc.CheckLanguage(ctx, data.LanguageCode)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.addWord - check language: %v", err))
		return
	}

	wordUUID, err := h.wordSvc.AddWord(ctx, &data)
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

	var data dto.GetWordRq

	if err := tools.CheckBody(w, r, &data); err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWord - check body: %v", err))
		return
	}

	ctx := r.Context()

	if err := h.langSvc.CheckLanguage(ctx, data.LanguageCode); err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWord - check language: %v", err))
		return
	}

	wordID, err := h.wordSvc.GetWord(ctx, &data)
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

	var data dto.RandomWordRq

	if err := tools.CheckBody(w, r, &data); err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getRandomWord - check body: %v", err))
		return
	}

	ctx := r.Context()

	if err := h.langSvc.CheckLanguage(ctx, data.LanguageCode); err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getRandomWord - check language: %v", err))
		return
	}

	word, err := h.wordSvc.GetRandomWord(ctx, &data)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getRandomWord: %v", err))
		return
	}

	randomWordRs := dto.RandomWordRs{
		Text:          word.Text,
		LanguageCode:  word.LanguageCode,
		Pronunciation: word.Pronunciation,
	}

	b, err := json.Marshal(&randomWordRs)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getRandomWord - marshal: %v", err))
		return
	}
	_, _ = w.Write(b)
}

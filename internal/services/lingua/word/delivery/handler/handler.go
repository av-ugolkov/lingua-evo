package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	entityLanguage "lingua-evo/internal/services/lingua/language"
	serviceLang "lingua-evo/internal/services/lingua/language/service"
	serviceWord "lingua-evo/internal/services/lingua/word/service"
	"lingua-evo/pkg/files"
	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/middleware"

	"github.com/google/uuid"
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
	AddWordRq struct {
		Text          string `json:"text"`
		LanguageCode  string `json:"language_code"`
		Pronunciation string `json:"pronunciation,omitempty"`
	}

	GetWordRq struct {
		Text         string `json:"text"`
		LanguageCode string `json:"language_code"`
	}
	GetWordIDRq struct {
		ID uuid.UUID `json:"id"`
	}

	RandomWordRq struct {
		LanguageCode string `json:"language_code"`
	}

	RandomWordRs struct {
		Text          string `json:"text"`
		LanguageCode  string `json:"language_code"`
		Pronunciation string `json:"pronunciation,omitempty"`
	}

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
	r.HandleFunc(addWord, middleware.Auth(h.addWord)).Methods(http.MethodPost)
	r.HandleFunc(getWord, h.getWord).Methods(http.MethodPost)
	r.HandleFunc(getRandomWord, h.getRandomWord).Methods(http.MethodPost)
}

func (h *Handler) openPage(w http.ResponseWriter, r *http.Request) {
	t, err := files.ParseFiles(addWordPage)
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
	defer func() {
		_ = r.Body.Close()
	}()

	handler := handler.NewHandler(w, r)
	var data AddWordRq
	if err := handler.CheckBody(&data); err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.addWord - check body: %v", err))
		return
	}

	ctx := r.Context()

	err := h.langSvc.CheckLanguage(ctx, data.LanguageCode)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.addWord - check language: %v", err))
		return
	}

	wordUUID, err := h.wordSvc.AddWord(ctx, data.Text, data.LanguageCode, data.Pronunciation)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, err)
		return
	}
	handler.SendData(http.StatusOK, []byte(wordUUID.String()))
}

func (h *Handler) getWord(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	handler := handler.NewHandler(w, r)
	var data GetWordRq
	if err := handler.CheckBody(&data); err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWord - check body: %v", err))
		return
	}

	ctx := r.Context()

	if err := h.langSvc.CheckLanguage(ctx, data.LanguageCode); err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWord - check language: %v", err))
		return
	}

	wordID, err := h.wordSvc.GetWord(ctx, data.Text, data.LanguageCode)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWord: %v", err))
		return
	}
	handler.SendData(http.StatusOK, []byte(wordID.String()))
}

func (h *Handler) getRandomWord(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	handler := handler.NewHandler(w, r)
	var data RandomWordRq
	if err := handler.CheckBody(&data); err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getRandomWord - check body: %v", err))
		return
	}

	ctx := r.Context()

	if err := h.langSvc.CheckLanguage(ctx, data.LanguageCode); err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getRandomWord - check language: %v", err))
		return
	}

	word, err := h.wordSvc.GetRandomWord(ctx, data.LanguageCode)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getRandomWord: %v", err))
		return
	}

	randomWordRs := RandomWordRs{
		Text:          word.Text,
		LanguageCode:  word.LanguageCode,
		Pronunciation: word.Pronunciation,
	}

	b, err := json.Marshal(&randomWordRs)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getRandomWord - marshal: %v", err))
		return
	}
	handler.SendData(http.StatusOK, b)
}

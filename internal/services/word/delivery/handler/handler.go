package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	serviceLang "lingua-evo/internal/services/language/service"
	serviceWord "lingua-evo/internal/services/word/service"
	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/http/handler/common"
	"lingua-evo/pkg/middleware"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	addWord       = "/word/add"
	getWord       = "/word/get"
	getRandomWord = "/word/get_random"
)

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
	r.HandleFunc(addWord, middleware.Auth(h.addWord)).Methods(http.MethodPost)
	r.HandleFunc(getWord, h.getWord).Methods(http.MethodPost)
	r.HandleFunc(getRandomWord, h.getRandomWord).Methods(http.MethodPost)
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
	handler.SetContentType(common.ContentTypeJSON)
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
	handler.SetContentType(common.ContentTypeJSON)
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
	handler.SetContentType(common.ContentTypeJSON)
	handler.SendData(http.StatusOK, b)
}

package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/pkg/http/exchange"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/middleware"
	langSvc "github.com/av-ugolkov/lingua-evo/internal/services/language"
	wordSvc "github.com/av-ugolkov/lingua-evo/internal/services/word"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	addWord       = "/word/add"
	getWord       = "/word/get"
	getRandomWord = "/word/random"
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

	WordRs struct {
		ID uuid.UUID `json:"id"`
	}

	RandomWordRs struct {
		Text          string `json:"text"`
		LanguageCode  string `json:"language_code"`
		Pronunciation string `json:"pronunciation,omitempty"`
	}

	Handler struct {
		wordSvc *wordSvc.Service
		langSvc *langSvc.Service
	}
)

func Create(r *mux.Router, wordSvc *wordSvc.Service, langSvc *langSvc.Service) {
	h := newHandler(wordSvc, langSvc)
	h.register(r)
}

func newHandler(wordSvc *wordSvc.Service, langSvc *langSvc.Service) *Handler {
	return &Handler{
		wordSvc: wordSvc,
		langSvc: langSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(addWord, middleware.Auth(h.addWord)).Methods(http.MethodPost)
	r.HandleFunc(getWord, h.getWord).Methods(http.MethodPost)
	r.HandleFunc(getRandomWord, h.getRandomWord).Methods(http.MethodGet)
}

func (h *Handler) addWord(ctx context.Context, ex *exchange.Exchanger) {
	var data AddWordRq
	if err := ex.CheckBody(&data); err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.addWord - check body: %v", err))
		return
	}

	err := h.langSvc.CheckLanguage(ctx, data.LanguageCode)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.addWord - check language: %v", err))
		return
	}

	wordID, err := h.wordSvc.AddWord(ctx, uuid.New(), data.Text, data.LanguageCode, data.Pronunciation)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, err)
		return
	}

	wordRs := &WordRs{
		ID: wordID,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, wordRs)
}

func (h *Handler) getWord(w http.ResponseWriter, r *http.Request) {
	ex := exchange.NewExchanger(w, r)
	var data GetWordRq
	if err := ex.CheckBody(&data); err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWord - check body: %v", err))
		return
	}

	ctx := r.Context()

	if err := h.langSvc.CheckLanguage(ctx, data.LanguageCode); err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWord - check language: %v", err))
		return
	}

	wordID, err := h.wordSvc.GetWordByValue(ctx, data.Text, data.LanguageCode)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWord: %v", err))
		return
	}

	wordRs := &WordRs{
		ID: wordID,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, wordRs)
}

func (h *Handler) getRandomWord(w http.ResponseWriter, r *http.Request) {
	ex := exchange.NewExchanger(w, r)
	ctx := ex.Context()

	languageCode, err := ex.QueryParamString("language_code")
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getRandomWord - query language: %v", err))
		return
	}

	if err := h.langSvc.CheckLanguage(ctx, languageCode); err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getRandomWord - check language: %v", err))
		return
	}

	word, err := h.wordSvc.GetRandomWord(ctx, languageCode)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getRandomWord: %v", err))
		return
	}

	randomWordRs := &RandomWordRs{
		Text:          word.Text,
		LanguageCode:  word.LanguageCode,
		Pronunciation: word.Pronunciation,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, randomWordRs)
}

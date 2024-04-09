package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/http/exchange"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/middleware"
	dictionarySvc "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	"github.com/av-ugolkov/lingua-evo/internal/services/dictionary/model"

	"github.com/gorilla/mux"
)

type Handler struct {
	dictSvc *dictionarySvc.Service
}

func Create(r *mux.Router, dictSvc *dictionarySvc.Service) {
	h := newHandler(dictSvc)
	h.register(r)
}

func newHandler(dictSvc *dictionarySvc.Service) *Handler {
	return &Handler{
		dictSvc: dictSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(delivery.DictionaryWord, middleware.Auth(h.addWord)).Methods(http.MethodPost)
	r.HandleFunc(delivery.DictionaryWord, h.getWord).Methods(http.MethodGet)
	r.HandleFunc(delivery.GetRandomWord, h.getRandomWord).Methods(http.MethodGet)
}

func (h *Handler) addWord(ctx context.Context, ex *exchange.Exchanger) {
	var data model.WordRq
	if err := ex.CheckBody(&data); err != nil {
		ex.SendError(http.StatusBadRequest, fmt.Errorf("dictionary.delivery.Handler.addWord - check body: %v", err))
		return
	}

	wordID, err := h.dictSvc.AddWord(ctx, data)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.addWord: %v", err))
		return
	}

	wordRs := &model.WordIDRs{
		ID: wordID,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, wordRs)
}

func (h *Handler) getWord(w http.ResponseWriter, r *http.Request) {
	ex := exchange.NewExchanger(w, r)
	ctx := ex.Context()

	text, err := ex.QueryParamString("text")
	if err != nil {
		ex.SendError(http.StatusBadRequest, fmt.Errorf("dictionary.delivery.Handler.getWord - check body: %v", err))
		return
	}

	langCode, err := ex.QueryParamString("lang_code")
	if err != nil {
		ex.SendError(http.StatusBadRequest, fmt.Errorf("dictionary.delivery.Handler.getWord - check body: %v", err))
		return
	}

	wordID, err := h.dictSvc.GetWordByText(ctx, text, langCode)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getWord: %v", err))
		return
	}

	wordRs := &model.WordIDRs{
		ID: wordID,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, wordRs)
}

func (h *Handler) getRandomWord(w http.ResponseWriter, r *http.Request) {
	ex := exchange.NewExchanger(w, r)
	ctx := ex.Context()

	langCode, err := ex.QueryParamString("lang_code")
	if err != nil {
		slog.Warn(fmt.Sprintf("dictionary.delivery.Handler.getRandomWord - get lang_code: %v", err))
	}

	word, err := h.dictSvc.GetRandomWord(ctx, langCode)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getRandomWord: %v", err))
		return
	}

	randomWordRs := &model.WordRs{
		Text:          word.Text,
		LangCode:      word.LangCode,
		Pronunciation: word.Pronunciation,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, randomWordRs)
}

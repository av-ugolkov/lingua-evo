package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/av-ugolkov/lingua-evo/internal/pkg/http/exchange"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	"github.com/av-ugolkov/lingua-evo/runtime"
)

const (
	dictionaryOp     = "/account/dictionary"
	getAllDictionary = "/account/dictionaries"
)

const (
	ParamsName = "name"
)

type (
	DictionaryRq struct {
		ID   uuid.UUID   `json:"id"`
		Name string      `json:"name"`
		Tags []uuid.UUID `json:"tags"`
	}

	DictionaryIDRs struct {
		ID   uuid.UUID   `json:"id"`
		Tags []uuid.UUID `json:"tags"`
	}

	DictionaryRs struct {
		ID     uuid.UUID `json:"id"`
		UserID uuid.UUID `json:"user_id"`
		Name   string    `json:"name"`
		Tags   []string  `json:"tags"`
	}

	Handler struct {
		dictionarySvc *dictionary.Service
	}
)

func Create(r *mux.Router, dictionarySvc *dictionary.Service) {
	h := newHandler(dictionarySvc)
	h.register(r)
}

func newHandler(dictionarySvc *dictionary.Service) *Handler {
	return &Handler{
		dictionarySvc: dictionarySvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(dictionaryOp, middleware.Auth(h.addDictionary)).Methods(http.MethodPost)
	r.HandleFunc(dictionaryOp, middleware.Auth(h.deleteDictionary)).Methods(http.MethodDelete)
	r.HandleFunc(dictionaryOp, middleware.Auth(h.getDictionary)).Methods(http.MethodGet)
	r.HandleFunc(dictionaryOp, middleware.Auth(h.renameDictionary)).Methods(http.MethodPut)
	r.HandleFunc(getAllDictionary, middleware.Auth(h.getDictionaries)).Methods(http.MethodGet)
}

func (h *Handler) addDictionary(ctx context.Context, ex *exchange.Exchanger) {
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ex.SendError(http.StatusUnauthorized, fmt.Errorf("dictionary.delivery.Handler.addDictionary - unauthorized: %v", err))
		return
	}

	name, err := ex.QueryParamString(ParamsName)
	if err != nil {
		ex.SendError(http.StatusUnauthorized, fmt.Errorf("dictionary.delivery.Handler.addDictionary - get query [name]: %v", err))
		return
	}

	dictID, err := h.dictionarySvc.AddDictionary(ctx, userID, uuid.New(), name)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.addDictionary: %v", err))
	}

	dictRs := &DictionaryRs{
		ID:     dictID,
		UserID: userID,
		Name:   name,
		Tags:   []string{},
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, dictRs)
}

func (h *Handler) deleteDictionary(ctx context.Context, ex *exchange.Exchanger) {
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ex.SendError(http.StatusUnauthorized, fmt.Errorf("dictionary.delivery.Handler.deleteDictionary - unauthorized: %v", err))
		return
	}

	name, err := ex.QueryParamString(ParamsName)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.deleteDictionary - get query [name]: %v", err))
		return
	}

	err = h.dictionarySvc.DeleteDictionary(ctx, userID, name)
	switch {
	case errors.Is(err, dictionary.ErrDictionaryNotFound):
		ex.SendError(http.StatusNotFound, fmt.Errorf("dictionary.delivery.Handler.deleteDictionary: %v", err))
		return
	case err != nil:
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.deleteDictionary: %v", err))
		return
	}

	ex.SendEmptyData(http.StatusOK)
}

func (h *Handler) getDictionary(ctx context.Context, ex *exchange.Exchanger) {
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ex.SendError(http.StatusUnauthorized, fmt.Errorf("dictionary.delivery.Handler.getDictionary - unauthorized: %v", err))
		return
	}

	name, err := ex.QueryParamString(ParamsName)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getDictionary - get query [name]: %v", err))
		return
	}

	id, tags, err := h.dictionarySvc.GetDictionary(ctx, userID, name)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getDictionary: %v", err))
		return
	}
	if id == uuid.Nil {
		ex.SendError(http.StatusNotFound, fmt.Errorf("dictionary.delivery.Handler.getDictionary - dictionary not found: %v", err))
		return
	}

	dictRs := &DictionaryIDRs{
		ID:   id,
		Tags: tags,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, dictRs)
}

func (h *Handler) getDictionaries(ctx context.Context, ex *exchange.Exchanger) {
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ex.SendError(http.StatusUnauthorized, fmt.Errorf("dictionary.delivery.Handler.getAllDictionary - unauthorized: %v", err))
		return
	}

	dictionaries, err := h.dictionarySvc.GetDictionaries(ctx, userID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getAllDictionary: %v", err))
	}

	dictionariesRs := make([]DictionaryRs, 0, len(dictionaries))
	for _, dict := range dictionaries {
		dictionariesRs = append(dictionariesRs, DictionaryRs{
			ID:     dict.ID,
			UserID: dict.UserID,
			Name:   dict.Name,
			Tags:   dict.Tags,
		})
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, dictionariesRs)
}

func (h *Handler) renameDictionary(ctx context.Context, ex *exchange.Exchanger) {
	name, err := ex.QueryParamString(ParamsName)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.renameDictionary - get query [name]: %v", err))
		return
	}

	var dictID DictionaryRq
	err = ex.CheckBody(&dictID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.renameDictionary - get body: %v", err))
		return
	}

	err = h.dictionarySvc.RenameDictionary(ctx, dictID.ID, name)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.renameDictionary: %v", err))
		return
	}

	ex.SendEmptyData(http.StatusOK)
}

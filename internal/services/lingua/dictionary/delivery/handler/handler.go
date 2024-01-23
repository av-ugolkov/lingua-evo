package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	entity "lingua-evo/internal/services/lingua/dictionary"
	"lingua-evo/internal/services/lingua/dictionary/service"
	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/http/handler/common"
	"lingua-evo/pkg/middleware"
	"lingua-evo/runtime"
)

const (
	dictionaryOp     = "/account/dictionary"
	getAllDictionary = "/account/dictionaries"
)

const (
	ParamsName     = "name"
	ParamsCapacity = "cap"
)

type (
	DictionaryRs struct {
		DictionaryID uuid.UUID `json:"dictionary_id"`
	}

	DictionaryIDRs struct {
		ID    uuid.UUID `json:"dictionary_id"`
		Words []string  `json:"words"`
	}

	Handler struct {
		dictionarySvc *service.DictionarySvc
	}
)

func Create(r *mux.Router, dictionarySvc *service.DictionarySvc) {
	handler := newHandler(dictionarySvc)
	handler.register(r)
}

func newHandler(dictionarySvc *service.DictionarySvc) *Handler {
	return &Handler{
		dictionarySvc: dictionarySvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(dictionaryOp, middleware.Auth(h.addDictionary)).Methods(http.MethodPost)
	r.HandleFunc(dictionaryOp, middleware.Auth(h.deleteDictionary)).Methods(http.MethodDelete)
	r.HandleFunc(dictionaryOp, middleware.Auth(h.getDictionary)).Methods(http.MethodGet)
	r.HandleFunc(getAllDictionary, middleware.Auth(h.getDictionaries)).Methods(http.MethodGet)
}

func (h *Handler) addDictionary(w http.ResponseWriter, r *http.Request) {
	handler := handler.NewHandler(w, r)
	ctx := r.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		handler.SendError(http.StatusUnauthorized, fmt.Errorf("dictionary.delivery.Handler.addDictionary - unauthorized: %v", err))
		return
	}

	name, err := handler.QueryParamString(ParamsName)
	if err != nil {
		handler.SendError(http.StatusUnauthorized, fmt.Errorf("dictionary.delivery.Handler.addDictionary - get query [name]: %v", err))
		return
	}

	dictID, err := h.dictionarySvc.AddDictionary(ctx, userID, name)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.addDictionary: %v", err))
	}

	b, err := json.Marshal(DictionaryRs{dictID})
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.addDictionary - marshal: %v", err))
		return
	}

	handler.SetContentType(common.ContentTypeJSON)
	handler.SendData(http.StatusOK, b)
}

func (h *Handler) deleteDictionary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	handler := handler.NewHandler(w, r)
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		handler.SendError(http.StatusUnauthorized, fmt.Errorf("dictionary.delivery.Handler.deleteDictionary - unauthorized: %v", err))
		return
	}

	name, err := handler.QueryParamString(ParamsName)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.deleteDictionary - get query [name]: %v", err))
		return
	}

	err = h.dictionarySvc.DeleteDictionary(ctx, userID, name)
	switch {
	case errors.Is(err, entity.ErrDictionaryNotFound):
		handler.SendError(http.StatusNotFound, fmt.Errorf("dictionary.delivery.Handler.deleteDictionary: %v", err))
		return
	case err != nil:
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.deleteDictionary: %v", err))
		return
	}

	handler.SendEmptyData(http.StatusOK)
}

func (h *Handler) getDictionary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	handler := handler.NewHandler(w, r)
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		handler.SendError(http.StatusUnauthorized, fmt.Errorf("dictionary.delivery.Handler.getDictionary - unauthorized: %v", err))
		return
	}

	name, err := handler.QueryParamString(ParamsName)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getDictionary - get query [name]: %v", err))
		return
	}

	capacity, err := handler.QueryParamInt(ParamsCapacity)
	if err != nil {
		slog.Warn(fmt.Sprintf("dictionary.delivery.Handler.getDictionary - get query [cap]: %v", err))
	}

	id, words, err := h.dictionarySvc.GetDictionary(ctx, userID, name, capacity)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getDictionary: %v", err))
		return
	}
	if id == uuid.Nil {
		handler.SendError(http.StatusNotFound, fmt.Errorf("dictionary.delivery.Handler.getDictionary - dictionary not found: %v", err))
		return
	}

	dictID := DictionaryIDRs{
		ID:    id,
		Words: words,
	}
	b, err := json.Marshal(&dictID)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getDictionary - marshal: %v", err))
		return
	}

	handler.SetContentType(common.ContentTypeJSON)
	handler.SendData(http.StatusOK, b)
}

func (h *Handler) getDictionaries(w http.ResponseWriter, r *http.Request) {
	handler := handler.NewHandler(w, r)
	userID, err := runtime.UserIDFromContext(r.Context())
	if err != nil {
		handler.SendError(http.StatusUnauthorized, fmt.Errorf("dictionary.delivery.Handler.getAllDictionary - unauthorized: %v", err))
		return
	}

	ctx := r.Context()
	dictionaries, err := h.dictionarySvc.GetDictionaries(ctx, userID)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getAllDictionary: %v", err))
	}

	b, err := json.Marshal(dictionaries)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getAllDictionary - marshal: %v", err))
		return
	}

	handler.SetContentType(common.ContentTypeJSON)
	handler.SendData(http.StatusOK, b)
}

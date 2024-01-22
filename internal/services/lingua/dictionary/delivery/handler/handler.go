package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"lingua-evo/internal/services/lingua/dictionary/service"
	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/http/handler/common"
	"lingua-evo/pkg/middleware"
	"lingua-evo/runtime"
)

const (
	addDictionary    = "/account/dictionary/add"
	deleteDictionary = "/account/dictionary/delete"
	getDictionary    = "/account/dictionary"
	getAllDictionary = "/account/dictionaries"
)

type (
	DictionaryRq struct {
		Name     string `json:"name"`
		Capacity int    `json:"cap"`
	}

	DictionariesRs struct {
		Dictionaries []Dictionary
	}

	Dictionary struct {
		ID   uuid.UUID
		Name string
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
	r.HandleFunc(addDictionary, middleware.Auth(h.addDictionary)).Methods(http.MethodPost)
	r.HandleFunc(deleteDictionary, middleware.Auth(h.deleteDictionary)).Methods(http.MethodDelete)
	r.HandleFunc(getDictionary, middleware.Auth(h.getDictionary)).Methods(http.MethodGet)
	r.HandleFunc(getAllDictionary, middleware.Auth(h.getDictionaries)).Methods(http.MethodGet)
}

func (h *Handler) addDictionary(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()
	handler := handler.NewHandler(w, r)
	ctx := r.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		handler.SendError(http.StatusUnauthorized, fmt.Errorf("dictionary.delivery.Handler.addDictionary - unauthorized: %v", err))
		return

	}

	var data DictionaryRq
	if err := handler.CheckBody(&data); err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.addDictionary - check body: %v", err))
		return
	}

	dictID, err := h.dictionarySvc.AddDictionary(ctx, userID, data.Name)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.addDictionary: %v", err))
	}

	handler.SetContentType(common.ContentTypeJSON)
	handler.SendData(http.StatusOK, []byte(dictID.String()))
}

func (h *Handler) deleteDictionary(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()
	ctx := r.Context()
	handler := handler.NewHandler(w, r)
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		handler.SendError(http.StatusUnauthorized, fmt.Errorf("dictionary.delivery.Handler.deleteDictionary - unauthorized: %v", err))
		return
	}

	var data DictionaryRq
	if err := handler.CheckBody(&data); err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.deleteDictionary - check body: %v", err))
		return
	}

	err = h.dictionarySvc.DeleteDictionary(ctx, userID, data.Name)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.deleteDictionary: %v", err))
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

	name, err := handler.QueryParams("name")
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getDictionary - get query [name]: %v", err))
		return
	}

	cap, err := handler.QueryParams("cap")
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getDictionary - get query [cap]: %v", err))
		return
	}
	capacity, err := strconv.Atoi(cap)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getDictionary - parse str to int: %v", err))
		return
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

package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"lingua-evo/internal/services/lingua/dictionary/service"
	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/middleware"
	"lingua-evo/runtime"
)

const (
	addDictionary    = "/dictionary/add"
	deleteDictionary = "/dictionary/delete"
	getDictionary    = "/dictionary/get"
	getAllDictionary = "/dictionary/get_all"
)

type (
	DictionaryRq struct {
		Name string `json:"name"`
	}

	DictionariesRs struct {
		Dictionaries []Dictionary
	}

	Dictionary struct {
		ID   uuid.UUID
		Name string
	}

	DictionaryIDRs struct {
		ID uuid.UUID `json:"dictionary_id"`
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
	r.HandleFunc(getDictionary, middleware.Auth(h.getDictionary)).Methods(http.MethodPost)
	r.HandleFunc(getAllDictionary, middleware.Auth(h.getAllDictionary)).Methods(http.MethodPost)
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

	handler.SendData([]byte(dictID.String()))
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
	handler.SendData([]byte("done"))
}

func (h *Handler) getDictionary(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()
	ctx := r.Context()
	handler := handler.NewHandler(w, r)
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		handler.SendError(http.StatusUnauthorized, fmt.Errorf("dictionary.delivery.Handler.getDictionary - unauthorized: %v", err))
		return
	}

	var data DictionaryRq
	if err := handler.CheckBody(&data); err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getDictionary - check body: %v", err))
		return
	}

	id, err := h.dictionarySvc.GetDictionary(ctx, userID, data.Name)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getDictionary: %v", err))
		return
	}

	dictID := DictionaryIDRs{
		ID: id,
	}
	b, err := json.Marshal(&dictID)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getDictionary - marshal: %v", err))
		return
	}
	handler.SendData(b)
}

func (h *Handler) getAllDictionary(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()
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

	//TODO нужно возвращать сериализованные данные
	slog.Info(fmt.Sprintf("count dictionaries: %d", len(dictionaries)))

	handler.SendData([]byte("done"))
}

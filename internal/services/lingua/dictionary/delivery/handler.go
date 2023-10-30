package delivery

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"

	"lingua-evo/internal/services/lingua/dictionary/dto"
	"lingua-evo/internal/services/lingua/dictionary/service"
	"lingua-evo/runtime"

	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/middleware"
)

const (
	addDictionary    = "/dictionary/add"
	deleteDictionary = "/dictionary/delete"
	getDictionary    = "/dictionary/get"
	getAllDictionary = "/dictionary/get_all"
)

type (
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
	r.Handle(addDictionary, middleware.AuthMiddleware(http.HandlerFunc(h.addDictionary))).Methods(http.MethodPost)
	r.Handle(deleteDictionary, middleware.AuthMiddleware(http.HandlerFunc(h.deleteDictionary))).Methods(http.MethodDelete)
	r.Handle(getDictionary, middleware.AuthMiddleware(http.HandlerFunc(h.getDictionary))).Methods(http.MethodPost)
	r.Handle(getAllDictionary, middleware.AuthMiddleware(http.HandlerFunc(h.getAllDictionary))).Methods(http.MethodPost)
}

func (h *Handler) addDictionary(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()
	ctx := r.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		handler.SendError(w, http.StatusUnauthorized, fmt.Errorf("dictionary.delivery.Handler.addDictionary - unauthorized: %v", err))
		return

	}

	var data dto.DictionaryRq
	if err := handler.CheckBody(w, r, &data); err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.addDictionary - check body: %v", err))
		return
	}

	dictID, err := h.dictionarySvc.AddDictionary(ctx, userID, &data)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.addDictionary: %v", err))
	}

	_, _ = w.Write([]byte(dictID.String()))
}

func (h *Handler) deleteDictionary(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()
	ctx := r.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		handler.SendError(w, http.StatusUnauthorized, fmt.Errorf("dictionary.delivery.Handler.deleteDictionary - unauthorized: %v", err))
		return
	}

	var data dto.DictionaryRq
	if err := handler.CheckBody(w, r, &data); err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.deleteDictionary - check body: %v", err))
		return
	}

	err = h.dictionarySvc.DeleteDictionary(ctx, userID, &data)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.deleteDictionary: %v", err))
	}

	_, _ = w.Write([]byte("done"))
}

func (h *Handler) getDictionary(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()
	ctx := r.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		handler.SendError(w, http.StatusUnauthorized, fmt.Errorf("dictionary.delivery.Handler.getDictionary - unauthorized: %v", err))
		return
	}

	var data dto.DictionaryRq
	if err := handler.CheckBody(w, r, &data); err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getDictionary - check body: %v", err))
		return
	}

	id, err := h.dictionarySvc.GetDictionary(ctx, userID, data.Name)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getDictionary: %v", err))
		return
	}

	dictID := dto.DictionaryIDRs{
		ID: id,
	}
	b, err := json.Marshal(&dictID)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getDictionary - marshal: %v", err))
		return
	}
	_, _ = w.Write(b)
}

func (h *Handler) getAllDictionary(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()
	userID, err := runtime.UserIDFromContext(r.Context())
	if err != nil {
		handler.SendError(w, http.StatusUnauthorized, fmt.Errorf("dictionary.delivery.Handler.getAllDictionary - unauthorized: %v", err))
		return
	}

	ctx := r.Context()
	dictionaries, err := h.dictionarySvc.GetDictionaries(ctx, userID)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getAllDictionary: %v", err))
	}

	//TODO нужно возвращать сериализованные данные
	slog.Info(fmt.Sprintf("count dictionaries: %d", len(dictionaries)))

	_, _ = w.Write([]byte("done"))
}

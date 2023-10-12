package delivery

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"lingua-evo/internal/services/dictionary/dto"
	"lingua-evo/internal/services/dictionary/entity"
	"lingua-evo/internal/tools"
)

const (
	addDictionary    = "/dictionary/add"
	deleteDictionary = "/dictionary/delete"
	getDictionary    = "/dictionary/get"
	getAllDictionary = "/dictionary/get_all"
)

type (
	dictionarySvc interface {
		AddDictionary(ctx context.Context, data dto.DictionaryRq) (uuid.UUID, error)
		DeleteDictionary(ctx context.Context, userID uuid.UUID, name string) error
		GetDictionary(ctx context.Context, userID uuid.UUID, name string) (uuid.UUID, error)
		GetDictionaries(ctx context.Context, userID uuid.UUID) ([]*entity.Dictionary, error)
	}

	Handler struct {
		dictionarySvc dictionarySvc
	}
)

func Create(r *mux.Router, dictionarySvc dictionarySvc) {
	handler := newHandler(dictionarySvc)
	handler.register(r)
}

func newHandler(dictionarySvc dictionarySvc) *Handler {
	return &Handler{
		dictionarySvc: dictionarySvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(addDictionary, h.addDictionary).Methods(http.MethodPost)
	r.HandleFunc(deleteDictionary, h.deleteDictionary).Methods(http.MethodPost)
	r.HandleFunc(getDictionary, h.getDictionary).Methods(http.MethodPost)
	r.HandleFunc(getAllDictionary, h.getAllDictionary).Methods(http.MethodPost)
}

func (h *Handler) addDictionary(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	var data dto.DictionaryRq
	err := tools.CheckBody(w, r, &data)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.addDictionary - check body: %v", err))
		return
	}

	ctx := r.Context()
	dictID, err := h.dictionarySvc.AddDictionary(ctx, data)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.addDictionary: %v", err))
	}

	_, _ = w.Write([]byte(dictID.String()))
}

func (h *Handler) deleteDictionary(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	var data dto.DictionaryRq
	err := tools.CheckBody(w, r, &data)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.deleteDictionary - check body: %v", err))
		return
	}

	ctx := r.Context()
	err = h.dictionarySvc.DeleteDictionary(ctx, data.UserID, data.Name)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.deleteDictionary: %v", err))
	}

	_, _ = w.Write([]byte("done"))
}

func (h *Handler) getDictionary(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	var data dto.DictionaryRq
	err := tools.CheckBody(w, r, &data)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getDictionary - check body: %v", err))
		return
	}

	ctx := r.Context()
	dictID, err := h.dictionarySvc.GetDictionary(ctx, data.UserID, data.Name)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getDictionary: %v", err))
	}

	_, _ = w.Write([]byte(dictID.String()))
}

func (h *Handler) getAllDictionary(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	var data uuid.UUID
	err := tools.CheckBody(w, r, &data)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getAllDictionary - check body: %v", err))
		return
	}

	ctx := r.Context()
	dictionaries, err := h.dictionarySvc.GetDictionaries(ctx, data)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("dictionary.delivery.Handler.getAllDictionary: %v", err))
	}

	//TODO нужно возвращать сериализованные данные
	slog.Info(fmt.Sprintf("count dictionaries: %d", len(dictionaries)))

	_, _ = w.Write([]byte("done"))
}

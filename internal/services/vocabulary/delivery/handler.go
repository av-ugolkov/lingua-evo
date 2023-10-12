package delivery

import (
	"context"
	"fmt"
	"net/http"

	"lingua-evo/internal/services/vocabulary/dto"
	"lingua-evo/internal/tools"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	addVocabulary    = "/vocabulary/add"
	deleteVocabulary = "/vocabulary/delete"
	getVocabulary    = "/vocabulary/get"
	getAllVocabulary = "/vocabulary/get_all"
)

type (
	vocabularySvc interface {
		AddWordInVocabulary(ctx context.Context, vocab *dto.AddWordRq) (uuid.UUID, error)
	}

	Handler struct {
		vocabularySvc vocabularySvc
	}
)

func Create(r *mux.Router, vocabularySvc vocabularySvc) {
	handler := newHandler(vocabularySvc)
	handler.register(r)
}

func newHandler(vocabularySvc vocabularySvc) *Handler {
	return &Handler{
		vocabularySvc: vocabularySvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(addVocabulary, h.addWord).Methods(http.MethodPost)
}

func (h *Handler) addWord(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	ctx := r.Context()
	var data dto.AddWordRq

	err := tools.CheckBody(w, r, &data)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.addWord - check body: %v", err))
		return
	}

	vocabID, err := h.vocabularySvc.AddWordInVocabulary(ctx, &data)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.addWord: %v", err))
		return
	}

	_, _ = w.Write([]byte(vocabID.String()))
}

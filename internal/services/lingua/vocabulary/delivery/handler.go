package delivery

import (
	"fmt"
	"net/http"

	"lingua-evo/internal/services/lingua/vocabulary/dto"
	"lingua-evo/internal/services/lingua/vocabulary/service"
	"lingua-evo/internal/tools"

	"github.com/gorilla/mux"
)

const (
	addVocabulary    = "/vocabulary/add"
	deleteVocabulary = "/vocabulary/delete"
	getVocabulary    = "/vocabulary/get"
	getAllVocabulary = "/vocabulary/get_all"
)

type (
	Handler struct {
		vocabularySvc *service.VocabularySvc
	}
)

func Create(r *mux.Router, vocabularySvc *service.VocabularySvc) {
	handler := newHandler(vocabularySvc)
	handler.register(r)
}

func newHandler(vocabularySvc *service.VocabularySvc) *Handler {
	return &Handler{
		vocabularySvc: vocabularySvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(addVocabulary, h.addWord).Methods(http.MethodPost)
	r.HandleFunc(deleteVocabulary, h.deleteWord).Methods(http.MethodDelete)
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

	err = h.vocabularySvc.AddWordInVocabulary(ctx, &data)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.addWord: %v", err))
		return
	}

	_, _ = w.Write([]byte("done"))
}

func (h *Handler) deleteWord(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	ctx := r.Context()
	var data dto.RemoveWordRq
	err := tools.CheckBody(w, r, &data)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.deleteWord - check body: %v", err))
		return
	}

	err = h.vocabularySvc.DeleteWordFromVocabulary(ctx, &data)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.deleteWord: %v", err))
		return
	}

	_, _ = w.Write([]byte("done"))
}

package delivery

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"lingua-evo/internal/services/word/entity"
)

const (
	addDictionary    = "/add_dictionary"
	removeDictionary = "/remove_dictionary"
)

type (
	dictionarySvc interface {
		AddDictionary(ctx context.Context, word *entity.Word) (uuid.UUID, error)
		RemoveDictionary(ctx context.Context, lang string) (*entity.Word, error)
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
	r.HandleFunc(removeDictionary, h.removeDictionary).Methods(http.MethodPost)
}

func (h *Handler) addDictionary(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) removeDictionary(w http.ResponseWriter, r *http.Request) {

}

package handler

import (
	"fmt"
	"net/http"

	entity "lingua-evo/internal/services/lingua/vocabulary"
	"lingua-evo/internal/services/lingua/vocabulary/service"
	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/middleware"

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
	AddWordRq struct {
		DictionaryID  uuid.UUID               `json:"dictionary_id"`
		NativeWord    entity.VocabularyWord   `json:"native_word"`
		TanslateWords []entity.VocabularyWord `json:"translate_words"`
		Examples      []string                `json:"examples"`
		Tags          []string                `json:"tags"`
	}

	RemoveWordRq struct {
		DictionaryID uuid.UUID `json:"dictionary_id"`
		NativeWordID uuid.UUID `json:"native_word_id"`
	}

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
	r.HandleFunc(addVocabulary, middleware.Auth(h.addWord)).Methods(http.MethodPost)
	r.HandleFunc(deleteVocabulary, middleware.Auth(h.deleteWord)).Methods(http.MethodDelete)
}

func (h *Handler) addWord(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	ctx := r.Context()
	handler := handler.NewHandler(w, r)

	var data AddWordRq
	err := handler.CheckBody(&data)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.addWord - check body: %v", err))
		return
	}

	err = h.vocabularySvc.AddWordInVocabulary(ctx, data.DictionaryID, data.NativeWord, data.TanslateWords, data.Examples, data.Tags)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.addWord: %v", err))
		return
	}
	handler.SendData([]byte("done"))
}

func (h *Handler) deleteWord(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	ctx := r.Context()
	handler := handler.NewHandler(w, r)

	var data RemoveWordRq
	err := handler.CheckBody(&data)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.deleteWord - check body: %v", err))
		return
	}

	err = h.vocabularySvc.DeleteWordFromVocabulary(ctx, data.DictionaryID, data.NativeWordID)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.deleteWord: %v", err))
		return
	}
	handler.SendData([]byte("done"))
}

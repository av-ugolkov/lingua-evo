package handler

import (
	"fmt"
	"net/http"

	entity "lingua-evo/internal/services/vocabulary"
	"lingua-evo/internal/services/vocabulary/service"
	"lingua-evo/pkg/http/exchange"
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
		DictionaryID  uuid.UUID     `json:"dictionary_id"`
		NativeWord    entity.Word   `json:"native_word"`
		TanslateWords []entity.Word `json:"translate_words"`
		Examples      []string      `json:"examples"`
		Tags          []string      `json:"tags"`
	}

	RemoveWordRq struct {
		DictionaryID uuid.UUID `json:"dictionary_id"`
		NativeWordID uuid.UUID `json:"native_word_id"`
	}

	VocabularyWordsRs struct {
		DictionaryId   uuid.UUID   `json:"dictionary_id"`
		NativeWord     uuid.UUID   `json:"native_word_id"`
		TranslateWords []uuid.UUID `json:"translate_words_id"`
		Examples       []uuid.UUID `json:"examples_id"`
		Tags           []uuid.UUID `json:"tags_id"`
	}

	Handler struct {
		vocabularySvc *service.VocabularySvc
	}
)

func Create(r *mux.Router, vocabularySvc *service.VocabularySvc) {
	h := newHandler(vocabularySvc)
	h.register(r)
}

func newHandler(vocabularySvc *service.VocabularySvc) *Handler {
	return &Handler{
		vocabularySvc: vocabularySvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(addVocabulary, middleware.Auth(h.addWord)).Methods(http.MethodPost)
	r.HandleFunc(deleteVocabulary, middleware.Auth(h.deleteWord)).Methods(http.MethodDelete)
	r.HandleFunc(getVocabulary, middleware.Auth(h.getWord)).Methods(http.MethodGet)
	r.HandleFunc(getAllVocabulary, middleware.Auth(h.getWords)).Methods(http.MethodGet)
}

func (h *Handler) addWord(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ex := exchange.NewExchanger(w, r)

	var data AddWordRq
	err := ex.CheckBody(&data)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.addWord - check body: %v", err))
		return
	}

	word, err := h.vocabularySvc.AddWordInVocabulary(ctx, data.DictionaryID, data.NativeWord, data.TanslateWords, data.Examples, data.Tags)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.addWord: %v", err))
		return
	}

	wordRs := &VocabularyWordsRs{
		DictionaryId:   word.DictionaryId,
		NativeWord:     word.NativeWord,
		TranslateWords: word.TranslateWords,
		Examples:       word.Examples,
		Tags:           word.Tags,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusCreated, wordRs)
}

func (h *Handler) deleteWord(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	ctx := r.Context()
	ex := exchange.NewExchanger(w, r)

	var data RemoveWordRq
	err := ex.CheckBody(&data)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.deleteWord - check body: %v", err))
		return
	}

	err = h.vocabularySvc.DeleteWordFromVocabulary(ctx, data.DictionaryID, data.NativeWordID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.deleteWord: %v", err))
		return
	}
	ex.SendEmptyData(http.StatusOK)
}

func (h *Handler) getWord(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getWords(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ex := exchange.NewExchanger(w, r)

	dictID, err := ex.QueryParamString("dictionary_id")
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.getWords - get query: %w", err))
		return
	}

	did, err := uuid.Parse(dictID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.getWords - parse query: %w", err))
		return
	}

	words, err := h.vocabularySvc.GetWords(ctx, did)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.getWords: %w", err))
		return
	}

	wordsRs := make([]VocabularyWordsRs, 0, len(words))
	for _, word := range words {
		wordRs := VocabularyWordsRs{
			DictionaryId:   word.DictionaryId,
			NativeWord:     word.NativeWord,
			TranslateWords: word.TranslateWords,
			Examples:       word.Examples,
			Tags:           word.Tags,
		}

		wordsRs = append(wordsRs, wordRs)
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, wordsRs)
}

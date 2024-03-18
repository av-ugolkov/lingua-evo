package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/pkg/http/exchange"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	ParamDictID = "dict_id"
	ParamWordID = "word_id"
	ParamName   = "name"
	ParamLimit  = "limit"
)

const (
	vocabularyWordUrl = "/vocabulary"
	getSeveralWords   = "/vocabulary/get_several_words"
	getWords          = "/vocabulary/all"
)

type (
	AddWordRq struct {
		DictionaryID  uuid.UUID       `json:"dictionary_id"`
		NativeWord    vocabulary.Word `json:"native_word"`
		TanslateWords []string        `json:"translate_words,omitempty"`
		Examples      []string        `json:"examples,omitempty"`
		Tags          []string        `json:"tags,omitempty"`
	}

	UpdateWordRq struct {
		OldWordID uuid.UUID `json:"word_id"`
		AddWordRq
	}

	RemoveWordRq struct {
		DictionaryID uuid.UUID `json:"dictionary_id"`
		NativeWordID uuid.UUID `json:"native_word_id"`
	}

	VocabularyWordsRs struct {
		WordID         *uuid.UUID      `json:"word_id,omitempty"`
		NativeWord     vocabulary.Word `json:"native"`
		TranslateWords []string        `json:"translate_words,omitempty"`
		Examples       []string        `json:"examples,omitempty"`
		Tags           []string        `json:"tags,omitempty"`
	}

	Handler struct {
		vocabularySvc *vocabulary.Service
	}
)

func Create(r *mux.Router, vocabularySvc *vocabulary.Service) {
	h := newHandler(vocabularySvc)
	h.register(r)
}

func newHandler(vocabularySvc *vocabulary.Service) *Handler {
	return &Handler{
		vocabularySvc: vocabularySvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(vocabularyWordUrl, middleware.Auth(h.getWord)).Methods(http.MethodGet)
	r.HandleFunc(vocabularyWordUrl, middleware.Auth(h.addWord)).Methods(http.MethodPut)
	r.HandleFunc(vocabularyWordUrl, middleware.Auth(h.deleteWord)).Methods(http.MethodDelete)
	r.HandleFunc(vocabularyWordUrl, middleware.Auth(h.updateWord)).Methods(http.MethodPatch)
	r.HandleFunc(getSeveralWords, middleware.Auth(h.getSeveralWords)).Methods(http.MethodGet)
	r.HandleFunc(getWords, middleware.Auth(h.getWords)).Methods(http.MethodGet)
}

func (h *Handler) addWord(ctx context.Context, ex *exchange.Exchanger) {
	var data AddWordRq
	err := ex.CheckBody(&data)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.addWord - check body: %v", err))
		return
	}

	translateWords := make([]vocabulary.Word, 0, len(data.TanslateWords))
	for _, word := range data.TanslateWords {
		translateWords = append(translateWords, vocabulary.Word{Text: word})
	}

	word, err := h.vocabularySvc.AddWord(ctx, data.DictionaryID, data.NativeWord, translateWords, data.Examples, data.Tags)
	if err != nil {
		switch {
		case errors.Is(err, vocabulary.ErrDuplicate):
			ex.SendError(http.StatusConflict, fmt.Errorf("vocabulary.delivery.Handler.addWord: %v", err))
			return
		default:
			ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.addWord: %v", err))
			return
		}
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusCreated, word)
}

func (h *Handler) deleteWord(ctx context.Context, ex *exchange.Exchanger) {
	var data RemoveWordRq
	err := ex.CheckBody(&data)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.deleteWord - check body: %v", err))
		return
	}

	err = h.vocabularySvc.DeleteWord(ctx, data.DictionaryID, data.NativeWordID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.deleteWord: %v", err))
		return
	}
	ex.SendEmptyData(http.StatusOK)
}

func (h *Handler) updateWord(ctx context.Context, ex *exchange.Exchanger) {
	var data UpdateWordRq
	err := ex.CheckBody(&data)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.updateWord - check body: %v", err))
		return
	}

	translateWords := make([]vocabulary.Word, 0, len(data.TanslateWords))
	for _, word := range data.TanslateWords {
		translateWords = append(translateWords, vocabulary.Word{Text: word})
	}

	word, err := h.vocabularySvc.UpdateWord(ctx, data.DictionaryID, data.OldWordID, data.NativeWord, translateWords, data.Examples, data.Tags)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.updateWord: %v", err))
		return
	}

	wordRs := &VocabularyWordsRs{
		NativeWord:     word.NativeWord,
		TranslateWords: word.TranslateWords,
		Examples:       word.Examples,
		Tags:           word.Tags,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusCreated, wordRs)
}

func (h *Handler) getSeveralWords(ctx context.Context, ex *exchange.Exchanger) {
	dictID, err := ex.QueryParamUUID(ParamDictID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.getSeveralWords - get dict id: %w", err))
		return
	}

	limit, err := ex.QueryParamInt(ParamLimit)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.getSeveralWords - get limit: %w", err))
		return
	}
	words, err := h.vocabularySvc.GetRandomWords(ctx, dictID, limit)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.getSeveralWords: %w", err))
		return
	}

	wordsRs := make([]VocabularyWordsRs, 0, len(words))
	for _, word := range words {
		wordRs := VocabularyWordsRs{
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

func (h *Handler) getWord(ctx context.Context, ex *exchange.Exchanger) {
	dictID, err := ex.QueryParamUUID(ParamDictID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.getWords - get dict id: %w", err))
		return
	}

	wordID, err := ex.QueryParamUUID(ParamWordID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.getWords - get word id: %w", err))
		return
	}
	word, err := h.vocabularySvc.GetWord(ctx, dictID, wordID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.getWords: %w", err))
		return
	}

	wordRs := VocabularyWordsRs{
		WordID:         &word.Id,
		NativeWord:     word.NativeWord,
		TranslateWords: word.TranslateWords,
		Examples:       word.Examples,
		Tags:           word.Tags,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, wordRs)
}

func (h *Handler) getWords(ctx context.Context, ex *exchange.Exchanger) {
	dictID, err := ex.QueryParamUUID(ParamDictID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.getWords - get dict id: %w", err))
		return
	}

	words, err := h.vocabularySvc.GetWords(ctx, dictID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.getWords: %w", err))
		return
	}

	wordsRs := make([]VocabularyWordsRs, 0, len(words))
	for _, word := range words {
		wordRs := VocabularyWordsRs{
			WordID:         &word.Id,
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

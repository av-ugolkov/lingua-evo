package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/http/exchange"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/services/word"
	"github.com/av-ugolkov/lingua-evo/internal/services/word/model"

	"github.com/gorilla/mux"
)

const (
	ParamDictID = "dict_id"
	ParamWordID = "word_id"
	ParamName   = "name"
	ParamLimit  = "limit"
)

type Handler struct {
	wordSvc *word.Service
}

func Create(r *mux.Router, wordSvc *word.Service) {
	h := newHandler(wordSvc)
	h.register(r)
}

func newHandler(wordSvc *word.Service) *Handler {
	return &Handler{
		wordSvc: wordSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(delivery.VocabularyWord, middleware.Auth(h.getWord)).Methods(http.MethodGet)
	r.HandleFunc(delivery.VocabularyWord, middleware.Auth(h.addWord)).Methods(http.MethodPost)
	r.HandleFunc(delivery.VocabularyWord, middleware.Auth(h.deleteWord)).Methods(http.MethodDelete)
	r.HandleFunc(delivery.VocabularyWord, middleware.Auth(h.updateWord)).Methods(http.MethodPatch)
	r.HandleFunc(delivery.VocabularySeveralWords, middleware.Auth(h.getSeveralWords)).Methods(http.MethodGet)
	r.HandleFunc(delivery.VocabularyWords, middleware.Auth(h.getWords)).Methods(http.MethodGet)
}

func (h *Handler) addWord(ctx context.Context, ex *exchange.Exchanger) {
	var data model.VocabWordRq
	err := ex.CheckBody(&data)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.addWord - check body: %v", err))
		return
	}

	vocabWord, err := h.wordSvc.AddWord(ctx, data)
	if err != nil {
		switch {
		case errors.Is(err, word.ErrDuplicate):
			ex.SendError(http.StatusConflict, fmt.Errorf("word.delivery.Handler.addWord: %v", err))
			return
		default:
			ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.addWord: %v", err))
			return
		}
	}

	wordRs := model.VocabWordsRs{
		WordID:         vocabWord.Id,
		NativeWord:     vocabWord.NativeWord,
		TranslateWords: vocabWord.TranslateWords,
		Examples:       vocabWord.Examples,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusCreated, wordRs)
}

func (h *Handler) deleteWord(ctx context.Context, ex *exchange.Exchanger) {
	var data model.RemoveVocabWordRq
	err := ex.CheckBody(&data)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.deleteWord - check body: %v", err))
		return
	}

	err = h.wordSvc.DeleteWord(ctx, data.VocabID, data.WordID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.deleteWord: %v", err))
		return
	}
	ex.SendEmptyData(http.StatusOK)
}

func (h *Handler) updateWord(ctx context.Context, ex *exchange.Exchanger) {
	var data model.VocabWordRq
	err := ex.CheckBody(&data)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.updateWord - check body: %v", err))
		return
	}

	translateWords := make([]model.Word, 0, len(data.TanslateWords))
	for _, tr := range data.TanslateWords {
		translateWords = append(translateWords, model.Word{Text: tr})
	}

	word, err := h.wordSvc.UpdateWord(ctx, data.VocabID, data.WordID, data.NativeWord, translateWords, data.Examples)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.updateWord: %v", err))
		return
	}

	wordRs := &model.VocabWordsRs{
		NativeWord:     word.NativeWord,
		TranslateWords: word.TranslateWords,
		Examples:       word.Examples,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusCreated, wordRs)
}

func (h *Handler) getSeveralWords(ctx context.Context, ex *exchange.Exchanger) {
	dictID, err := ex.QueryParamUUID(ParamDictID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getSeveralWords - get dict id: %w", err))
		return
	}

	limit, err := ex.QueryParamInt(ParamLimit)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getSeveralWords - get limit: %w", err))
		return
	}
	words, err := h.wordSvc.GetRandomWords(ctx, dictID, limit)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getSeveralWords: %w", err))
		return
	}

	wordsRs := make([]model.VocabWordsRs, 0, len(words))
	for _, word := range words {
		wordRs := model.VocabWordsRs{
			NativeWord:     word.NativeWord,
			TranslateWords: word.TranslateWords,
			Examples:       word.Examples,
		}

		wordsRs = append(wordsRs, wordRs)
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, wordsRs)
}

func (h *Handler) getWord(ctx context.Context, ex *exchange.Exchanger) {
	dictID, err := ex.QueryParamUUID(ParamDictID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWords - get dict id: %w", err))
		return
	}

	wordID, err := ex.QueryParamUUID(ParamWordID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWords - get word id: %w", err))
		return
	}
	word, err := h.wordSvc.GetWord(ctx, dictID, wordID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWords: %w", err))
		return
	}

	wordRs := model.VocabWordsRs{
		WordID:         word.Id,
		NativeWord:     word.NativeWord,
		TranslateWords: word.TranslateWords,
		Examples:       word.Examples,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, wordRs)
}

func (h *Handler) getWords(ctx context.Context, ex *exchange.Exchanger) {
	dictID, err := ex.QueryParamUUID(ParamDictID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWords - get dict id: %w", err))
		return
	}

	words, err := h.wordSvc.GetWords(ctx, dictID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWords: %w", err))
		return
	}

	wordsRs := make([]model.VocabWordsRs, 0, len(words))
	for _, word := range words {
		wordRs := model.VocabWordsRs{
			WordID:         word.Id,
			NativeWord:     word.NativeWord,
			TranslateWords: word.TranslateWords,
			Examples:       word.Examples,
		}

		wordsRs = append(wordsRs, wordRs)
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, wordsRs)
}

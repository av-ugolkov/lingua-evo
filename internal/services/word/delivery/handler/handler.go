package handler

import (
	"context"
	"errors"
	"fmt"
	entityDict "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/delivery"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/http/exchange"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/middleware"
	entityExample "github.com/av-ugolkov/lingua-evo/internal/services/example"
	"github.com/av-ugolkov/lingua-evo/internal/services/word"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/word"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	ParamVocabID = "vocab_id"
	ParamID      = "id"
	ParamLimit   = "limit"
)

type (
	VocabWord struct {
		ID            *uuid.UUID `json:"id,omitempty"`
		Text          string     `json:"text,omitempty"`
		Pronunciation string     `json:"pronunciation,omitempty"`
	}

	VocabWordRq struct {
		ID         *uuid.UUID `json:"id,omitempty"`
		VocabID    uuid.UUID  `json:"vocab_id"`
		Native     VocabWord  `json:"native"`
		Translates []string   `json:"translates,omitempty"`
		Examples   []string   `json:"examples,omitempty"`
	}

	RemoveVocabWordRq struct {
		VocabID uuid.UUID `json:"vocab_id"`
		WordID  uuid.UUID `json:"word_id"`
	}

	VocabWordRs struct {
		ID         *uuid.UUID `json:"id,omitempty"`
		Native     *VocabWord `json:"native,omitempty"`
		Translates []string   `json:"translates,omitempty"`
		Examples   []string   `json:"examples,omitempty"`
		Created    *time.Time `json:"created,omitempty"`
		Updated    *time.Time `json:"updated,omitempty"`
	}
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
	r.HandleFunc(delivery.VocabularyWordUpdate, middleware.Auth(h.updateWord)).Methods(http.MethodPost)
	r.HandleFunc(delivery.VocabularySeveralWords, middleware.Auth(h.getSeveralWords)).Methods(http.MethodGet)
	r.HandleFunc(delivery.VocabularyWords, middleware.Auth(h.getWords)).Methods(http.MethodGet)
}

func (h *Handler) addWord(ctx context.Context, ex *exchange.Exchanger) {
	var data VocabWordRq
	err := ex.CheckBody(&data)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.addWord - check body: %v", err))
		return
	}

	translateWords := make([]entityDict.DictWord, 0, len(data.Translates))
	for _, translateWord := range data.Translates {
		translateWords = append(translateWords, entityDict.DictWord{
			ID:        uuid.New(),
			Text:      translateWord,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
	}

	examples := make([]entityExample.Example, 0, len(data.Examples))
	for _, example := range data.Examples {
		examples = append(examples, entityExample.Example{
			ID:        uuid.New(),
			Text:      example,
			CreatedAt: time.Now().UTC(),
		})
	}

	vocabWord, err := h.wordSvc.AddWord(ctx, entity.VocabWordData{
		ID:      uuid.New(),
		VocabID: data.VocabID,
		Native: entityDict.DictWord{
			ID:            uuid.New(),
			Text:          data.Native.Text,
			Pronunciation: data.Native.Pronunciation,
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		},
		Translates: translateWords,
		Examples:   examples,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	})
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

	wordRs := VocabWordRs{
		ID: &vocabWord.ID,
		Native: &VocabWord{
			ID: &vocabWord.NativeID,
		},
		Created: &vocabWord.CreatedAt,
		Updated: &vocabWord.UpdatedAt,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusCreated, wordRs)
}

func (h *Handler) updateWord(ctx context.Context, ex *exchange.Exchanger) {
	var data VocabWordRq
	err := ex.CheckBody(&data)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.updateWord - check body: %v", err))
		return
	}

	translates := make([]entityDict.DictWord, 0, len(data.Translates))
	for _, tr := range data.Translates {
		translates = append(translates, entityDict.DictWord{
			ID:        uuid.New(),
			Text:      tr,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
	}

	examples := make([]entityExample.Example, 0, len(data.Examples))
	for _, example := range data.Examples {
		examples = append(examples, entityExample.Example{
			ID:        uuid.New(),
			Text:      example,
			CreatedAt: time.Now().UTC(),
		})
	}

	vocabWord, err := h.wordSvc.UpdateWord(ctx, entity.VocabWordData{
		ID:      *data.ID,
		VocabID: data.VocabID,
		Native: entityDict.DictWord{
			ID:            *data.Native.ID,
			Text:          data.Native.Text,
			Pronunciation: data.Native.Pronunciation,
			UpdatedAt:     time.Now().UTC(),
		},
		Translates: translates,
		Examples:   examples,
		UpdatedAt:  time.Now().UTC(),
	})
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.updateWord: %v", err))
		return
	}

	wordRs := &VocabWordRs{
		ID: &vocabWord.ID,
		Native: &VocabWord{
			ID: &vocabWord.NativeID,
		},
		Created: &vocabWord.CreatedAt,
		Updated: &vocabWord.UpdatedAt,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, wordRs)
}

func (h *Handler) deleteWord(ctx context.Context, ex *exchange.Exchanger) {
	var data RemoveVocabWordRq
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

func (h *Handler) getSeveralWords(ctx context.Context, ex *exchange.Exchanger) {
	vocabID, err := ex.QueryParamUUID(ParamVocabID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getSeveralWords - get vocab id: %w", err))
		return
	}

	limit, err := ex.QueryParamInt(ParamLimit)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getSeveralWords - get limit: %w", err))
		return
	}
	vocabWords, err := h.wordSvc.GetRandomWords(ctx, vocabID, limit)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getSeveralWords: %w", err))
		return
	}

	wordsRs := make([]VocabWordRs, 0, len(vocabWords))
	for _, vocabWord := range vocabWords {
		translates := make([]string, 0, len(vocabWord.Translates))
		for _, translate := range vocabWord.Translates {
			translates = append(translates, translate.Text)
		}

		wordRs := VocabWordRs{
			Native: &VocabWord{
				Text:          vocabWord.Native.Text,
				Pronunciation: vocabWord.Native.Pronunciation,
			},
			Translates: translates,
		}

		wordsRs = append(wordsRs, wordRs)
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, wordsRs)
}

func (h *Handler) getWord(ctx context.Context, ex *exchange.Exchanger) {
	wordID, err := ex.QueryParamUUID(ParamID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWords - get word id: %w", err))
		return
	}

	vocabWord, err := h.wordSvc.GetWord(ctx, wordID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWords: %w", err))
		return
	}

	translates := make([]string, 0, len(vocabWord.Translates))
	for _, translate := range vocabWord.Translates {
		translates = append(translates, translate.Text)
	}

	examples := make([]string, 0, len(vocabWord.Examples))
	for _, example := range vocabWord.Examples {
		examples = append(examples, example.Text)
	}

	wordRs := VocabWordRs{
		ID: &vocabWord.ID,
		Native: &VocabWord{
			ID:            &vocabWord.Native.ID,
			Text:          vocabWord.Native.Text,
			Pronunciation: vocabWord.Native.Pronunciation,
		},
		Translates: translates,
		Examples:   examples,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, wordRs)
}

func (h *Handler) getWords(ctx context.Context, ex *exchange.Exchanger) {
	vocabID, err := ex.QueryParamUUID(ParamVocabID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWords - get dict id: %w", err))
		return
	}

	vocabWords, err := h.wordSvc.GetWords(ctx, vocabID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("word.delivery.Handler.getWords: %w", err))
		return
	}

	wordsRs := make([]VocabWordRs, 0, len(vocabWords))
	for _, vocabWord := range vocabWords {
		translates := make([]string, 0, len(vocabWord.Translates))
		for _, translate := range vocabWord.Translates {
			translates = append(translates, translate.Text)
		}

		examples := make([]string, 0, len(vocabWord.Examples))
		for _, example := range vocabWord.Examples {
			examples = append(examples, example.Text)
		}

		wordRs := VocabWordRs{
			ID: &vocabWord.ID,
			Native: &VocabWord{
				ID:            &vocabWord.Native.ID,
				Text:          vocabWord.Native.Text,
				Pronunciation: vocabWord.Native.Pronunciation,
			},
			Translates: translates,
			Examples:   examples,
			Created:    &vocabWord.CreatedAt,
			Updated:    &vocabWord.UpdatedAt,
		}

		wordsRs = append(wordsRs, wordRs)
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, wordsRs)
}

package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/av-ugolkov/lingua-evo/internal/delivery"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/http/exchange"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/middleware"
	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	"github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime"
)

const (
	paramsVocabName = "name"
)

type (
	VocabularyRq struct {
		Name          string   `json:"name"`
		NativeLang    string   `json:"native_lang"`
		TranslateLang string   `json:"translate_lang"`
		Tags          []string `json:"tags"`
	}

	VocabularyIDRs struct {
		ID uuid.UUID `json:"id"`
	}

	VocabularyRs struct {
		ID            uuid.UUID `json:"id"`
		UserID        uuid.UUID `json:"user_id"`
		Name          string    `json:"name"`
		NativeLang    string    `json:"native_lang"`
		TranslateLang string    `json:"translate_lang"`
		Tags          []string  `json:"tags"`
	}
)

type Handler struct {
	vocabularySvc *vocabulary.Service
}

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
	r.HandleFunc(delivery.Vocabulary, middleware.Auth(h.addVocabulary)).Methods(http.MethodPost)
	r.HandleFunc(delivery.Vocabulary, middleware.Auth(h.deleteVocabulary)).Methods(http.MethodDelete)
	r.HandleFunc(delivery.Vocabulary, middleware.Auth(h.getVocabulary)).Methods(http.MethodGet)
	r.HandleFunc(delivery.Vocabulary, middleware.Auth(h.renameVocabulary)).Methods(http.MethodPut)
	r.HandleFunc(delivery.Vocabularies, middleware.Auth(h.getVocabularies)).Methods(http.MethodGet)
}

func (h *Handler) addVocabulary(ctx context.Context, ex *exchange.Exchanger) {
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ex.SendError(http.StatusUnauthorized, fmt.Errorf("vocabulary.delivery.Handler.addVocabulary - unauthorized: %v", err))
		return
	}

	var data VocabularyRq
	err = ex.CheckBody(&data)
	if err != nil {
		ex.SendError(http.StatusUnauthorized, fmt.Errorf("vocabulary.delivery.Handler.addVocabulary - check body: %v", err))
		return
	}

	tags := make([]entityTag.Tag, 0, len(data.Tags))
	for _, tag := range data.Tags {
		tags = append(tags, entityTag.Tag{
			ID:   uuid.New(),
			Text: tag,
		})
	}

	vocab, err := h.vocabularySvc.AddVocabulary(ctx, vocabulary.Vocabulary{
		ID:            uuid.New(),
		UserID:        userID,
		Name:          data.Name,
		NativeLang:    data.NativeLang,
		TranslateLang: data.TranslateLang,
		Tags:          tags,
	})
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.addVocabulary: %v", err))
	}

	vocabRs := &VocabularyRs{
		ID:            vocab.ID,
		UserID:        vocab.UserID,
		Name:          vocab.Name,
		NativeLang:    vocab.NativeLang,
		TranslateLang: vocab.TranslateLang,
		//Tags:          vocab.Tags,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, vocabRs)
}

func (h *Handler) deleteVocabulary(ctx context.Context, ex *exchange.Exchanger) {
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ex.SendError(http.StatusUnauthorized, fmt.Errorf("vocabulary.delivery.Handler.deleteVocabulary - unauthorized: %v", err))
		return
	}

	name, err := ex.QueryParamString(paramsVocabName)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.deleteVocabulary - get query [name]: %v", err))
		return
	}

	err = h.vocabularySvc.DeleteVocabulary(ctx, userID, name)
	switch {
	case errors.Is(err, vocabulary.ErrVocabularyNotFound):
		ex.SendError(http.StatusNotFound, fmt.Errorf("vocabulary.delivery.Handler.deleteVocabulary: %v", err))
		return
	case err != nil:
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.deleteVocabulary: %v", err))
		return
	}

	ex.SendEmptyData(http.StatusOK)
}

func (h *Handler) getVocabulary(ctx context.Context, ex *exchange.Exchanger) {
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ex.SendError(http.StatusUnauthorized, fmt.Errorf("vocabulary.delivery.Handler.getVocabulary - unauthorized: %v", err))
		return
	}

	name, err := ex.QueryParamString(paramsVocabName)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.getVocabulary - get query [name]: %v", err))
		return
	}

	vocab, err := h.vocabularySvc.GetVocabulary(ctx, userID, name)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.getVocabulary: %v", err))
		return
	}
	if vocab.ID == uuid.Nil {
		ex.SendError(http.StatusNotFound, fmt.Errorf("vocabulary.delivery.Handler.getVocabulary - vocabulary not found: %v", err))
		return
	}

	tags := make([]string, 0, len(vocab.Tags))
	for _, tag := range vocab.Tags {
		tags = append(tags, tag.Text)
	}

	vocabRs := &VocabularyRs{
		ID:            vocab.ID,
		UserID:        vocab.UserID,
		Name:          vocab.Name,
		NativeLang:    vocab.NativeLang,
		TranslateLang: vocab.TranslateLang,
		Tags:          tags,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, vocabRs)
}

func (h *Handler) getVocabularies(ctx context.Context, ex *exchange.Exchanger) {
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ex.SendError(http.StatusUnauthorized, fmt.Errorf("vocabulary.delivery.Handler.getVocabularies - unauthorized: %v", err))
		return
	}

	vocabularies, err := h.vocabularySvc.GetVocabularies(ctx, userID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.getVocabularies: %v", err))
	}

	vocabulariesRs := make([]VocabularyRs, 0, len(vocabularies))
	for _, vocab := range vocabularies {
		tags := make([]string, 0, len(vocab.Tags))
		for _, tag := range vocab.Tags {
			tags = append(tags, tag.Text)
		}

		vocabulariesRs = append(vocabulariesRs, VocabularyRs{
			ID:            vocab.ID,
			UserID:        vocab.UserID,
			Name:          vocab.Name,
			NativeLang:    vocab.NativeLang,
			TranslateLang: vocab.TranslateLang,
			Tags:          tags,
		})
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, vocabulariesRs)
}

func (h *Handler) renameVocabulary(ctx context.Context, ex *exchange.Exchanger) {
	name, err := ex.QueryParamString(paramsVocabName)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.renameVocabulary - get query [name]: %v", err))
		return
	}

	var vocab VocabularyIDRs
	err = ex.CheckBody(&vocab)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.renameVocabulary - get body: %v", err))
		return
	}

	err = h.vocabularySvc.RenameVocabulary(ctx, vocab.ID, name)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.renameVocabulary: %v", err))
		return
	}

	ex.SendEmptyData(http.StatusOK)
}

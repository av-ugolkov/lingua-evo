package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/http/exchange"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/middleware"
	tagSvc "github.com/av-ugolkov/lingua-evo/internal/services/tag"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	queryVocabID = "vocab_id"
)

type (
	TagRs struct {
		ID   uuid.UUID `json:"id"`
		Text string    `json:"text"`
	}
)

type Handler struct {
	tagSvc *tagSvc.Service
}

func Create(r *mux.Router, tagSvc *tagSvc.Service) {
	h := newHandler(tagSvc)
	h.register(r)
}

func newHandler(tagSvc *tagSvc.Service) *Handler {
	return &Handler{
		tagSvc: tagSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(delivery.VocabularyTags, middleware.Auth(h.getTags)).Methods(http.MethodGet)
}

func (h *Handler) getTags(ctx context.Context, ex *exchange.Exchanger) {
	vocabIDStr, err := ex.QueryParamString(queryVocabID)
	if err != nil {
		ex.SendError(http.StatusBadRequest, fmt.Errorf("tag.delivery.Handler.getTags - query param [%s]: %v", queryVocabID, err))
		return
	}

	vocabID, err := uuid.Parse(vocabIDStr)
	if err != nil {
		ex.SendError(http.StatusBadRequest, fmt.Errorf("tag.delivery.Handler.getTags - invalid id: %v", err))
		return
	}

	tags, err := h.tagSvc.GetTagsInVocabulary(ctx, vocabID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("tag.delivery.Handler.getTags: %v", err))
		return
	}

	tagsRs := make([]TagRs, 0, len(tags))
	for _, tag := range tags {
		tagsRs = append(tagsRs, TagRs{
			ID:   tag.ID,
			Text: tag.Text,
		})
	}

	ex.SendData(http.StatusOK, tagsRs)
}

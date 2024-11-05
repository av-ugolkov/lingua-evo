package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	tagSvc "github.com/av-ugolkov/lingua-evo/internal/services/tag"

	"github.com/google/uuid"
)

const (
	paramsVocabID = "vocab_id"
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

func Create(r *ginext.Engine, tagSvc *tagSvc.Service) {
	h := newHandler(tagSvc)

	r.GET(handler.VocabularyTags, middleware.Auth(h.getTags))
}

func newHandler(tagSvc *tagSvc.Service) *Handler {
	return &Handler{
		tagSvc: tagSvc,
	}
}

func (h *Handler) getTags(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	vocabID, err := c.GetQueryUUID(paramsVocabID)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("tag.delivery.Handler.getTags: %v", err)
	}

	tags, err := h.tagSvc.GetTagsInVocabulary(ctx, vocabID)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("tag.delivery.Handler.getTags: %v", err)
	}

	tagsRs := make([]TagRs, 0, len(tags))
	for _, tag := range tags {
		tagsRs = append(tagsRs, TagRs{
			ID:   tag.ID,
			Text: tag.Text,
		})
	}

	return http.StatusOK, tagsRs, nil
}

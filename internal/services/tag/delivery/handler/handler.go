package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/pkg/http/gin_extension"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/middleware"
	tagSvc "github.com/av-ugolkov/lingua-evo/internal/services/tag"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func Create(r *gin.Engine, tagSvc *tagSvc.Service) {
	h := newHandler(tagSvc)
	h.register(r)
}

func newHandler(tagSvc *tagSvc.Service) *Handler {
	return &Handler{
		tagSvc: tagSvc,
	}
}

func (h *Handler) register(r *gin.Engine) {
	r.GET(delivery.VocabularyTags, middleware.Auth(h.getTags))
}

func (h *Handler) getTags(c *gin.Context) {
	ctx := c.Request.Context()
	vocabIDStr, err := ginExt.GetQuery(c, queryVocabID)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("tag.delivery.Handler.getTags - query param [%s]: %v", queryVocabID, err))
		return
	}

	vocabID, err := uuid.Parse(vocabIDStr)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("tag.delivery.Handler.getTags - invalid id: %v", err))
		return
	}

	tags, err := h.tagSvc.GetTagsInVocabulary(ctx, vocabID)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("tag.delivery.Handler.getTags: %v", err))
		return
	}

	tagsRs := make([]TagRs, 0, len(tags))
	for _, tag := range tags {
		tagsRs = append(tagsRs, TagRs{
			ID:   tag.ID,
			Text: tag.Text,
		})
	}

	c.JSON(http.StatusOK, tagsRs)
}

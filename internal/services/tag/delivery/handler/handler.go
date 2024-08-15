package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/delivery/handler/gin"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	tagSvc "github.com/av-ugolkov/lingua-evo/internal/services/tag"

	"github.com/gin-gonic/gin"
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
	r.GET(handler.VocabularyTags, middleware.Auth(h.getTags))
}

func (h *Handler) getTags(c *gin.Context) {
	ctx := c.Request.Context()
	vocabIDStr, err := ginExt.GetQuery(c, paramsVocabID)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("tag.delivery.Handler.getTags - query param [%s]: %v", paramsVocabID, err))
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

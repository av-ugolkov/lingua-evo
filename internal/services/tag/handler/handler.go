package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	tagSvc "github.com/av-ugolkov/lingua-evo/internal/services/tag"

	"github.com/gofiber/fiber/v2"
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

func Create(r *fiber.App, tagSvc *tagSvc.Service) {
	h := newHandler(tagSvc)

	r.Get(handler.VocabularyTags, middleware.Auth(h.getTags))
}

func newHandler(tagSvc *tagSvc.Service) *Handler {
	return &Handler{
		tagSvc: tagSvc,
	}
}

func (h *Handler) getTags(c *fiber.Ctx) error {
	ctx := c.Context()

	vocabID, err := uuid.Parse(c.Query(paramsVocabID))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("tag.Handler.getTags: %v", err))
	}

	tags, err := h.tagSvc.GetTagsInVocabulary(ctx, vocabID)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("tag.Handler.getTags: %v", err))
	}

	tagsRs := make([]TagRs, 0, len(tags))
	for _, tag := range tags {
		tagsRs = append(tagsRs, TagRs{
			ID:   tag.ID,
			Text: tag.Text,
		})
	}

	return c.Status(http.StatusOK).JSON(tagsRs)
}

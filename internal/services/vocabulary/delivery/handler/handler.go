package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/av-ugolkov/lingua-evo/internal/delivery"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/pkg/http/gin_extension"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/middleware"
	vocabulary "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/service"
	"github.com/av-ugolkov/lingua-evo/runtime"
)

const (
	paramVocabName = "name"
	paramVocabID   = "id"
)

type (
	VocabularyRq struct {
		Name          string   `json:"name"`
		Access        int      `json:"access_id"`
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
		AccessID      int       `json:"access_id"`
		NativeLang    string    `json:"native_lang"`
		TranslateLang string    `json:"translate_lang"`
		Description   string    `json:"description"`
		Tags          []string  `json:"tags"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}

	VocabularyEditRq struct {
		ID     uuid.UUID `json:"id"`
		Name   string    `json:"name"`
		Access int       `json:"access_id"`
	}
)

type Handler struct {
	vocabularySvc *vocabulary.Service
}

func Create(r *gin.Engine, vocabularySvc *vocabulary.Service) {
	h := newHandler(vocabularySvc)
	h.register(r)
}

func newHandler(vocabularySvc *vocabulary.Service) *Handler {
	return &Handler{
		vocabularySvc: vocabularySvc,
	}
}

func (h *Handler) register(r *gin.Engine) {
	r.POST(delivery.UserVocabulary, middleware.Auth(h.userAddVocabulary))
	r.DELETE(delivery.UserVocabulary, middleware.Auth(h.userDeleteVocabulary))
	r.GET(delivery.UserVocabulary, middleware.Auth(h.userGetVocabulary))
	r.PUT(delivery.UserVocabulary, middleware.Auth(h.userEditVocabulary))
	r.GET(delivery.UserVocabularies, middleware.Auth(h.userGetVocabularies))
	r.GET(delivery.Vocabularies, middleware.OptionalAuth(h.getVocabularies))
}

func (h *Handler) getVocabularies(c *gin.Context) {
	ctx := c.Request.Context()
	userID, _ := runtime.UserIDFromContext(ctx)

	vocabularies, err := h.vocabularySvc.GetVocabularies(ctx, userID)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies: %v", err))
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
			AccessID:      vocab.Access,
			NativeLang:    vocab.NativeLang,
			TranslateLang: vocab.TranslateLang,
			Description:   vocab.Description,
			Tags:          tags,
			CreatedAt:     vocab.CreatedAt,
			UpdatedAt:     vocab.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, vocabulariesRs)
}

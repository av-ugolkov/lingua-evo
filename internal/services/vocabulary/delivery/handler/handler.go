package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/av-ugolkov/lingua-evo/internal/delivery"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/pkg/http/gin_extension"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/middleware"
	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	"github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime"
)

const (
	paramsVocabName = "name"
	paramsVocabID   = "id"
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
		Tags          []string  `json:"tags"`
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
	r.POST(delivery.Vocabulary, middleware.Auth(h.addVocabulary))
	r.DELETE(delivery.Vocabulary, middleware.Auth(h.deleteVocabulary))
	r.GET(delivery.Vocabulary, middleware.Auth(h.getVocabulary))
	r.PUT(delivery.Vocabulary, middleware.Auth(h.editVocabulary))
	r.GET(delivery.Vocabularies, middleware.Auth(h.getVocabularies))
}

func (h *Handler) addVocabulary(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ginExt.SendError(c, http.StatusUnauthorized,
			fmt.Errorf("vocabulary.delivery.Handler.addVocabulary - unauthorized: %v", err))
		return
	}

	var data VocabularyRq
	err = c.Bind(&data)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("vocabulary.delivery.Handler.addVocabulary - check body: %v", err))
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
		Access:        data.Access,
		NativeLang:    data.NativeLang,
		TranslateLang: data.TranslateLang,
		Tags:          tags,
	})
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.addVocabulary: %v", err))
		return
	}

	vocabRs := VocabularyRs{
		ID:            vocab.ID,
		UserID:        vocab.UserID,
		Name:          vocab.Name,
		NativeLang:    vocab.NativeLang,
		TranslateLang: vocab.TranslateLang,
		//Tags:          vocab.Tags,
	}

	c.JSON(http.StatusOK, vocabRs)
}

func (h *Handler) deleteVocabulary(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ginExt.SendError(c, http.StatusUnauthorized,
			fmt.Errorf("vocabulary.delivery.Handler.deleteVocabulary - unauthorized: %v", err))
		return
	}

	name, err := ginExt.GetQuery(c, paramsVocabName)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.deleteVocabulary - get query [name]: %v", err))
		return
	}

	err = h.vocabularySvc.DeleteVocabulary(ctx, userID, name)
	switch {
	case errors.Is(err, vocabulary.ErrVocabularyNotFound):
		ginExt.SendError(c, http.StatusNotFound,
			fmt.Errorf("vocabulary.delivery.Handler.deleteVocabulary: %v", err))
		return
	case err != nil:
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.deleteVocabulary: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) getVocabulary(c *gin.Context) {
	ctx := c.Request.Context()

	vocabID, err := ginExt.GetQueryUUID(c, paramsVocabID)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabulary - get query [name]: %v", err))
		return
	}

	vocab, err := h.vocabularySvc.GetVocabulary(ctx, vocabID)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabulary: %v", err))
		return
	}
	if vocab.ID == uuid.Nil {
		ginExt.SendError(c, http.StatusNotFound,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabulary - vocabulary not found: %v", err))
		return
	}

	tags := make([]string, 0, len(vocab.Tags))
	for _, tag := range vocab.Tags {
		tags = append(tags, tag.Text)
	}

	vocabRs := VocabularyRs{
		ID:            vocab.ID,
		UserID:        vocab.UserID,
		Name:          vocab.Name,
		NativeLang:    vocab.NativeLang,
		TranslateLang: vocab.TranslateLang,
		Tags:          tags,
	}

	c.JSON(http.StatusOK, vocabRs)
}

func (h *Handler) getVocabularies(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ginExt.SendError(c, http.StatusUnauthorized,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies - unauthorized: %v", err))
		return
	}

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
			Tags:          tags,
		})
	}

	c.JSON(http.StatusOK, vocabulariesRs)
}

func (h *Handler) editVocabulary(c *gin.Context) {
	ctx := c.Request.Context()

	var data VocabularyEditRq
	err := c.Bind(&data)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("vocabulary.delivery.Handler.editVocabulary - check body: %v", err))
		return
	}

	err = h.vocabularySvc.EditVocabulary(ctx, vocabulary.Vocabulary{
		ID:     data.ID,
		Name:   data.Name,
		Access: data.Access,
	})
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.editVocabulary: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

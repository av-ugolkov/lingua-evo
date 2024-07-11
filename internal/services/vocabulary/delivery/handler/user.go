package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	ginExt "github.com/av-ugolkov/lingua-evo/internal/pkg/http/gin_extension"
	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	"github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime"
)

func (h *Handler) userAddVocabulary(c *gin.Context) {
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

	vocab, err := h.vocabularySvc.UserAddVocabulary(ctx, vocabulary.Vocabulary{
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

func (h *Handler) userDeleteVocabulary(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ginExt.SendError(c, http.StatusUnauthorized,
			fmt.Errorf("vocabulary.delivery.Handler.deleteVocabulary - unauthorized: %v", err))
		return
	}

	name, err := ginExt.GetQuery(c, paramVocabName)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.deleteVocabulary - get query [name]: %v", err))
		return
	}

	err = h.vocabularySvc.UserDeleteVocabulary(ctx, userID, name)
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

func (h *Handler) userGetVocabulary(c *gin.Context) {
	ctx := c.Request.Context()

	vocabID, err := ginExt.GetQueryUUID(c, paramVocabID)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabulary - get query [name]: %v", err))
		return
	}

	vocab, err := h.vocabularySvc.UserGetVocabulary(ctx, vocabID)
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

func (h *Handler) userGetVocabularies(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ginExt.SendError(c, http.StatusUnauthorized,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies - unauthorized: %v", err))
		return
	}

	vocabularies, err := h.vocabularySvc.UserGetVocabularies(ctx, userID)
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

func (h *Handler) userEditVocabulary(c *gin.Context) {
	ctx := c.Request.Context()

	var data VocabularyEditRq
	err := c.Bind(&data)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("vocabulary.delivery.Handler.editVocabulary - check body: %v", err))
		return
	}

	err = h.vocabularySvc.UserEditVocabulary(ctx, vocabulary.Vocabulary{
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

package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	ginExt "github.com/av-ugolkov/lingua-evo/internal/delivery/handler/gin"
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

	vocab, err := h.vocabSvc.UserAddVocabulary(ctx, vocabulary.Vocab{
		UserID:        userID,
		Name:          data.Name,
		Access:        data.Access,
		NativeLang:    data.NativeLang,
		TranslateLang: data.TranslateLang,
		Description:   data.Description,
		Tags:          tags,
	})
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.addVocabulary: %v", err))
		return
	}

	vocabRs := VocabularyRs{
		ID: vocab.ID,
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

	name, err := ginExt.GetQuery(c, paramsVocabName)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.deleteVocabulary - get query [name]: %v", err))
		return
	}

	err = h.vocabSvc.UserDeleteVocabulary(ctx, userID, name)
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

func (h *Handler) userGetVocabularies(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ginExt.SendError(c, http.StatusUnauthorized,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies - unauthorized: %v", err))
		return
	}

	page, err := ginExt.GetQueryInt(c, paramsPage)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies - get query [page]: %v", err))
		return
	}

	itemsPerPage, err := ginExt.GetQueryInt(c, paramsPerPage)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies - get query [per_page]: %v", err))
		return
	}

	typeSort, err := ginExt.GetQueryInt(c, paramsSort)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies - get query [order]: %v", err))
		return
	}

	order, err := ginExt.GetQueryInt(c, paramsOrder)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies - get query [order]: %v", err))
		return
	}

	search, err := ginExt.GetQuery(c, paramsSearch)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies - get query [search]: %v", err))
		return
	}

	nativeLang, err := ginExt.GetQuery(c, paramsNativeLang)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies - get query [native_lang]: %v", err))
		return
	}

	translateLang, err := ginExt.GetQuery(c, paramsTranslateLang)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies - get query [translate_lang]: %v", err))
		return
	}

	vocabs, totalCount, err := h.vocabSvc.UserGetVocabularies(ctx, userID, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies: %v", err))
	}

	vocabulariesRs := make([]VocabularyRs, 0, len(vocabs))
	for _, vocab := range vocabs {
		vocabulariesRs = append(vocabulariesRs, VocabularyRs{
			ID:            vocab.ID,
			UserID:        vocab.UserID,
			Name:          vocab.Name,
			AccessID:      &vocab.Access,
			NativeLang:    vocab.NativeLang,
			TranslateLang: vocab.TranslateLang,
			CreatedAt:     vocab.CreatedAt,
			UserName:      vocab.UserName,
			WordsCount:    &vocab.WordsCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"vocabularies": vocabulariesRs,
		"total_count":  totalCount,
	})
}

func (h *Handler) userEditVocabulary(c *gin.Context) {
	ctx := c.Request.Context()

	var data VocabularyRq
	err := c.Bind(&data)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("vocabulary.delivery.Handler.editVocabulary - check body: %v", err))
		return
	}

	err = h.vocabSvc.UserEditVocabulary(ctx, vocabulary.Vocab{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		Access:      data.Access,
	})
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError, fmt.Errorf("vocabulary.delivery.Handler.editVocabulary: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

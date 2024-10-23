package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	"github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gin-gonic/gin"
)

func (h *Handler) userAddVocabulary(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			fmt.Errorf("vocabulary.delivery.Handler.addVocabulary: %v", err)
	}

	var data VocabularyRq
	err = c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("vocabulary.delivery.Handler.addVocabulary: %v", err)
	}

	tags := make([]entityTag.Tag, 0, len(data.Tags))
	for _, tag := range data.Tags {
		tags = append(tags, entityTag.Tag{
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
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.addVocabulary: %v", err)
	}

	return http.StatusOK, VocabularyRs{
		ID: vocab.ID,
	}, nil
}

func (h *Handler) userDeleteVocabulary(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			fmt.Errorf("vocabulary.delivery.Handler.deleteVocabulary: %v", err)
	}

	name, err := c.GetQuery(paramsVocabName)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.deleteVocabulary: %v", err)
	}

	err = h.vocabSvc.UserDeleteVocabulary(ctx, userID, name)
	switch {
	case errors.Is(err, vocabulary.ErrVocabularyNotFound):
		return http.StatusNotFound, nil,
			fmt.Errorf("vocabulary.delivery.Handler.deleteVocabulary: %v", err)
	case err != nil:
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.deleteVocabulary: %v", err)
	}

	return http.StatusOK, gin.H{}, nil
}

func (h *Handler) userGetVocabularies(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies: %v", err)
	}

	page, err := c.GetQueryInt(paramsPage)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies: %v", err)
	}

	itemsPerPage, err := c.GetQueryInt(paramsPerPage)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies: %v", err)
	}

	typeSort, err := c.GetQueryInt(paramsSort)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies: %v", err)
	}

	order, err := c.GetQueryInt(paramsOrder)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies: %v", err)
	}

	search, err := c.GetQuery(paramsSearch)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies: %v", err)
	}

	nativeLang, err := c.GetQuery(paramsNativeLang)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies: %v", err)
	}

	translateLang, err := c.GetQuery(paramsTranslateLang)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies: %v", err)
	}

	vocabs, totalCount, err := h.vocabSvc.UserGetVocabularies(ctx, userID, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies: %v", err)
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

	return http.StatusOK, gin.H{
		"vocabularies": vocabulariesRs,
		"total_count":  totalCount,
	}, nil
}

func (h *Handler) userEditVocabulary(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	var data VocabularyRq
	err := c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("vocabulary.delivery.Handler.editVocabulary: %v", err)
	}

	err = h.vocabSvc.UserEditVocabulary(ctx, vocabulary.Vocab{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		Access:      data.Access,
	})
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("vocabulary.delivery.Handler.editVocabulary: %v", err)
	}

	return http.StatusOK, gin.H{}, nil
}

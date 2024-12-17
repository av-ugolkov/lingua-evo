package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	"github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gofiber/fiber/v2"
)

const (
	ErrMsgNameTooLong        = "Vocabulary name should be less than 50 characters"
	ErrMsgDescriptionTooLong = "Description should be less than 250 characters"
	ErrMsgCountTags          = "Count of tags should be less than 5"
)

func (h *Handler) userAddVocabulary(c *fiber.Ctx) error {
	ctx := c.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized,
			fmt.Sprintf("vocabulary.Handler.addVocabulary: %v", err))
	}

	var data VocabularyRq
	err = c.BodyParser(&data)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, msgerr.ErrMsgBadRequest)
	}

	if len(data.Name) > 50 {
		return fiber.NewError(http.StatusBadRequest, ErrMsgNameTooLong)
	}

	if len(data.Description) > 250 {
		return fiber.NewError(http.StatusBadRequest, ErrMsgDescriptionTooLong)
	}

	vocab, err := h.vocabSvc.UserAddVocabulary(ctx, vocabulary.Vocab{
		UserID:        userID,
		Name:          data.Name,
		Access:        data.Access,
		NativeLang:    data.NativeLang,
		TranslateLang: data.TranslateLang,
		Description:   data.Description,
	})
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, msgerr.ErrMsgInternal)
	}

	return c.Status(http.StatusOK).JSON(VocabularyRs{
		ID:        &vocab.ID,
		UserID:    &vocab.UserID,
		CreatedAt: &vocab.CreatedAt,
		UpdatedAt: &vocab.UpdatedAt,
	})
}

func (h *Handler) userDeleteVocabulary(c *fiber.Ctx) error {
	ctx := c.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized,
			fmt.Sprintf("vocabulary.Handler.deleteVocabulary: %v", err))
	}

	name := c.Query(paramsVocabName)
	if name == runtime.EmptyString {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.deleteVocabulary: %v", err))
	}

	err = h.vocabSvc.UserDeleteVocabulary(ctx, userID, name)
	switch {
	case errors.Is(err, vocabulary.ErrVocabularyNotFound):
		return fiber.NewError(http.StatusNotFound,
			fmt.Sprintf("vocabulary.Handler.deleteVocabulary: %v", err))
	case err != nil:
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.deleteVocabulary: %v", err))
	}

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) userGetVocabularies(c *fiber.Ctx) error {
	ctx := c.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized, msgerr.ErrMsgUnauthorized)
	}

	page := c.QueryInt(paramsPage)
	if page == 0 {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.getVocabularies: %v", err))
	}

	itemsPerPage := c.QueryInt(paramsPerPage)
	if itemsPerPage == 0 {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.getVocabularies: %v", err))
	}

	typeSort := c.QueryInt(paramsSort, -1)
	if typeSort == -1 {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.getVocabularies: %v", err))
	}

	order := c.QueryInt(paramsOrder, -1)
	if order == -1 {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.getVocabularies: %v", err))
	}

	search := c.Query(paramsSearch)
	nativeLang := c.Query(paramsNativeLang)
	if nativeLang == runtime.EmptyString {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.getVocabularies: %v", err))
	}

	translateLang := c.Query(paramsTranslateLang)
	if translateLang == runtime.EmptyString {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.getVocabularies: %v", err))
	}

	vocabs, totalCount, err := h.vocabSvc.UserGetVocabularies(ctx, userID, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.getVocabularies: %v", err))
	}

	vocabulariesRs := make([]VocabularyRs, 0, len(vocabs))
	for _, vocab := range vocabs {
		vocabulariesRs = append(vocabulariesRs, VocabularyRs{
			ID:            &vocab.ID,
			UserID:        &vocab.UserID,
			Name:          vocab.Name,
			AccessID:      &vocab.Access,
			NativeLang:    vocab.NativeLang,
			TranslateLang: vocab.TranslateLang,
			CreatedAt:     &vocab.CreatedAt,
			UserName:      vocab.UserName,
			WordsCount:    &vocab.WordsCount,
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"vocabularies": vocabulariesRs,
		"total_count":  totalCount,
	})
}

func (h *Handler) userEditVocabulary(c *fiber.Ctx) error {
	ctx := c.Context()

	var data VocabularyRq
	err := c.BodyParser(&data)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, msgerr.ErrMsgBadRequest)
	}

	if len(data.Name) > 50 {
		return fiber.NewError(http.StatusBadRequest, ErrMsgNameTooLong)
	}

	if len(data.Description) > 250 {
		return fiber.NewError(http.StatusBadRequest, ErrMsgDescriptionTooLong)
	}

	err = h.vocabSvc.UserEditVocabulary(ctx, vocabulary.Vocab{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		Access:      data.Access,
	})
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.editVocabulary: %v", err))
	}

	return c.SendStatus(http.StatusOK)
}

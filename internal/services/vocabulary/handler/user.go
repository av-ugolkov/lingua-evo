package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/pkg/fext"
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
	userID, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	var data VocabularyRq
	err = c.BodyParser(&data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	if len(data.Name) > 50 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.addVocabulary: name is too long"),
			ErrMsgNameTooLong))
	}

	if len(data.Description) > 250 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.addVocabulary: description is too long"),
			ErrMsgDescriptionTooLong))
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
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgInternal))
	}

	return c.Status(http.StatusOK).JSON(fext.D(VocabularyRs{
		ID:        &vocab.ID,
		UserID:    &vocab.UserID,
		CreatedAt: &vocab.CreatedAt,
		UpdatedAt: &vocab.UpdatedAt,
	}))
}

func (h *Handler) userDeleteVocabulary(c *fiber.Ctx) error {
	ctx := c.Context()
	userID, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	name := c.Query(paramsVocabName)
	if name == runtime.EmptyString {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.deleteVocabulary: %v", err)))
	}

	err = h.vocabSvc.UserDeleteVocabulary(ctx, userID, name)
	switch {
	case errors.Is(err, vocabulary.ErrVocabularyNotFound):
		return c.Status(http.StatusNotFound).JSON(fext.E(err))
	case err != nil:
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) userGetVocabularies(c *fiber.Ctx) error {
	ctx := c.Context()
	userID, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	page := c.QueryInt(paramsPage)
	if page == 0 {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.getVocabularies: not found query [%s]", paramsPage)))
	}

	itemsPerPage := c.QueryInt(paramsPerPage)
	if itemsPerPage == 0 {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.getVocabularies: not found query [%s]", paramsPerPage)))
	}

	typeSort := c.QueryInt(paramsSort, -1)
	if typeSort == -1 {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.getVocabularies: not found query [%s]", paramsSort)))
	}

	order := c.QueryInt(paramsOrder, -1)
	if order == -1 {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.getVocabularies: not found query [%s]", paramsOrder)))
	}

	search := c.Query(paramsSearch)
	nativeLang := c.Query(paramsNativeLang)
	if nativeLang == runtime.EmptyString {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.getVocabularies: not found query [%s]", paramsNativeLang)))
	}

	translateLang := c.Query(paramsTranslateLang)
	if translateLang == runtime.EmptyString {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.getVocabularies: not found query [%s]", paramsTranslateLang)))
	}

	vocabs, totalCount, err := h.vocabSvc.UserGetVocabularies(ctx, userID, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
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

	return c.Status(http.StatusOK).JSON(fext.D(fiber.Map{
		"vocabularies": vocabulariesRs,
		"total_count":  totalCount,
	}))
}

func (h *Handler) userEditVocabulary(c *fiber.Ctx) error {
	ctx := c.Context()

	var data VocabularyRq
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	if len(data.Name) > 50 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.editVocabulary: name is too long"),
			ErrMsgNameTooLong))
	}

	if len(data.Description) > 250 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.editVocabulary: description is too long"),
			ErrMsgDescriptionTooLong))
	}

	err = h.vocabSvc.UserEditVocabulary(ctx, vocabulary.Vocab{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		Access:      data.Access,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	return c.SendStatus(http.StatusOK)
}

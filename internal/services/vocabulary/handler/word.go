package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/pkg/fext"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils/slices"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	ErrWordTooLong          = "Word or phrase length should be less than 100 characters"
	ErrPronunciationTooLong = "Pronunciation length should be less than 100 characters"
	ErrDefinitionTooLong    = "Definition length should be less than 100 characters"
	ErrCountTranslates      = "Count of translates should be less than 10"
	ErrCountExamples        = "Count of examples should be less than 5"
	ErrWordIsExists         = "This word is already exists"
	ErrExampleIsExists      = "This example is already exists"
)

type (
	VocabWord struct {
		ID            *uuid.UUID `json:"id,omitempty"`
		Text          string     `json:"text,omitempty"`
		Pronunciation string     `json:"pronunciation,omitempty"`
	}

	RemoveVocabWordRq struct {
		VocabID uuid.UUID `json:"vocab_id"`
		WordID  uuid.UUID `json:"word_id"`
	}

	VocabWordRs struct {
		ID         *uuid.UUID `json:"id,omitempty"`
		Native     *VocabWord `json:"native,omitempty"`
		Definition string     `json:"definition,omitempty"`
		Translates []string   `json:"translates,omitempty"`
		Examples   []string   `json:"examples,omitempty"`
		Created    int64      `json:"created,omitempty"`
		Updated    int64      `json:"updated,omitempty"`
	}
)

func (h *Handler) addWord(c *fiber.Ctx) error {
	ctx := c.Context()

	userID, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	var data struct {
		VocabID    uuid.UUID `json:"vocab_id"`
		Native     VocabWord `json:"native"`
		Definition string    `json:"definition,omitempty"`
		Translates []string  `json:"translates,omitempty"`
		Examples   []string  `json:"examples,omitempty"`
	}
	err = c.BodyParser(&data)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	if len(data.Native.Text) > 100 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.addWord: native word is too long"),
			ErrWordTooLong))
	}

	if len(data.Native.Pronunciation) > 100 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.addWord: pronunciation is too long"),
			ErrPronunciationTooLong))
	}

	if len(data.Definition) > 100 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.addWord: definition is too long"), ErrDefinitionTooLong))
	}

	if len(data.Translates) > 10 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.addWord: count of translates is too long"), ErrCountTranslates))
	}

	if slices.HasDuplicates(data.Translates) {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.addWord: word is already exists"), ErrWordIsExists))
	}

	translates := make([]entity.DictWord, 0, len(data.Translates))
	for _, translate := range data.Translates {
		translates = append(translates, entity.DictWord{
			Text: translate,
		})
	}

	if len(data.Examples) > 5 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.addWord: count of examples is too long"),
			ErrCountExamples))
	}
	examples := make([]entity.Example, 0, len(data.Examples))
	for _, example := range data.Examples {
		examples = append(examples, entity.Example{
			Text: example,
		})
	}

	vocabWord, err := h.vocabSvc.AddWord(ctx, userID, entity.VocabWordData{
		VocabID: data.VocabID,
		Native: entity.DictWord{
			Text:          data.Native.Text,
			Pronunciation: data.Native.Pronunciation,
		},
		Definition: data.Definition,
		Translates: translates,
		Examples:   examples,
	})
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrDuplicate):
			return c.Status(http.StatusConflict).JSON(fext.E(err, ErrWordIsExists))
		default:
			return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgInternal))
		}
	}

	wordRs := VocabWordRs{
		ID: &vocabWord.ID,
		Native: &VocabWord{
			ID: &vocabWord.NativeID,
		},
		Created: vocabWord.CreatedAt.UnixMilli(),
		Updated: vocabWord.UpdatedAt.UnixMilli(),
	}

	return c.Status(http.StatusCreated).JSON(fext.D(wordRs))
}

func (h *Handler) updateWordText(c *fiber.Ctx) error {
	ctx := c.Context()
	userID, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	var data struct {
		ID      uuid.UUID `json:"id"`
		VocabID uuid.UUID `json:"vocab_id"`
		Native  struct {
			ID   uuid.UUID `json:"id"`
			Text string    `json:"text"`
		}
	}
	err = c.BodyParser(&data)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	if len(data.Native.Text) > 100 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.updateWordText: native word is too long"),
			ErrWordTooLong))
	}

	vocabWord, err := h.vocabSvc.UpdateWordText(ctx, userID, entity.VocabWordData{
		ID:      data.ID,
		VocabID: data.VocabID,
		Native: entity.DictWord{
			ID:   data.Native.ID,
			Text: data.Native.Text,
		},
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	wordRs := &VocabWordRs{
		ID: &vocabWord.ID,
		Native: &VocabWord{
			ID: &vocabWord.NativeID,
		},
		Updated: vocabWord.UpdatedAt.UnixMilli(),
	}

	return c.Status(http.StatusOK).JSON(fext.D(wordRs))
}

func (h *Handler) updateWordPronunciation(c *fiber.Ctx) error {
	ctx := c.Context()
	userID, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	var data struct {
		ID      uuid.UUID `json:"id"`
		VocabID uuid.UUID `json:"vocab_id"`
		Native  struct {
			ID            uuid.UUID `json:"id"`
			Text          string    `json:"text"`
			Pronunciation string    `json:"pronunciation"`
		}
	}
	err = c.BodyParser(&data)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	if len(data.Native.Pronunciation) > 100 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.updateWordPronunciation: native pronunciation is too long"),
			ErrPronunciationTooLong))
	}

	vocabWord, err := h.vocabSvc.UpdateWordPronunciation(ctx, userID, entity.VocabWordData{
		ID:      data.ID,
		VocabID: data.VocabID,
		Native: entity.DictWord{
			ID:            data.Native.ID,
			Text:          data.Native.Text,
			Pronunciation: data.Native.Pronunciation,
		},
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	wordRs := &VocabWordRs{
		ID:      &vocabWord.ID,
		Updated: vocabWord.UpdatedAt.UnixMilli(),
	}

	return c.Status(http.StatusOK).JSON(fext.D(wordRs))
}

func (h *Handler) updateWordDefinition(c *fiber.Ctx) error {
	ctx := c.Context()
	userID, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	var data struct {
		ID      uuid.UUID `json:"id"`
		VocabID uuid.UUID `json:"vocab_id"`
		Native  struct {
			ID   uuid.UUID `json:"id"`
			Text string    `json:"text"`
		}
		Definition string `json:"definition"`
	}
	err = c.BodyParser(&data)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	if len(data.Definition) > 100 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.updateWordDefinition: definition is too long"),
			ErrDefinitionTooLong))
	}

	vocabWord, err := h.vocabSvc.UpdateWordDefinition(ctx, userID, entity.VocabWordData{
		ID:         data.ID,
		VocabID:    data.VocabID,
		Definition: data.Definition,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	wordRs := &VocabWordRs{
		ID:      &vocabWord.ID,
		Updated: vocabWord.UpdatedAt.UnixMilli(),
	}

	return c.Status(http.StatusOK).JSON(fext.D(wordRs))
}

func (h *Handler) updateWordTranslates(c *fiber.Ctx) error {
	ctx := c.Context()
	userID, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	var data struct {
		ID         uuid.UUID `json:"id"`
		VocabID    uuid.UUID `json:"vocab_id"`
		Translates []string  `json:"translates"`
	}
	err = c.BodyParser(&data)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	if len(data.Translates) > 10 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.updateWordTranslates: count of translates is too long"),
			ErrCountTranslates))
	}

	if slices.HasDuplicates(data.Translates) {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.updateWordTranslates: translates are not unique"),
			ErrWordIsExists))
	}

	translates := make([]entity.DictWord, 0, len(data.Translates))
	for _, translate := range data.Translates {
		translates = append(translates, entity.DictWord{
			Text: translate,
		})
	}

	vocabWord, err := h.vocabSvc.UpdateWordTranslates(ctx, userID, entity.VocabWordData{
		ID:         data.ID,
		VocabID:    data.VocabID,
		Translates: translates,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	wordRs := &VocabWordRs{
		ID:      &vocabWord.ID,
		Updated: vocabWord.UpdatedAt.UnixMilli(),
	}

	return c.Status(http.StatusOK).JSON(fext.D(wordRs))
}

func (h *Handler) updateWordExamples(c *fiber.Ctx) error {
	ctx := c.Context()
	userID, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	var data struct {
		ID       uuid.UUID `json:"id"`
		VocabID  uuid.UUID `json:"vocab_id"`
		Examples []string  `json:"examples"`
	}
	err = c.BodyParser(&data)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	if len(data.Examples) > 5 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.updateWordExamples: count of examples is too long"),
			ErrCountExamples))
	}

	if slices.HasDuplicates(data.Examples) {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("vocabulary.Handler.updateWordExamples: examples are not unique"),
			ErrExampleIsExists))
	}

	examples := make([]entity.Example, 0, len(data.Examples))
	for _, example := range data.Examples {
		examples = append(examples, entity.Example{
			Text: example,
		})
	}

	vocabWord, err := h.vocabSvc.UpdateWordExamples(ctx, userID, entity.VocabWordData{
		ID:       data.ID,
		VocabID:  data.VocabID,
		Examples: examples,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	wordRs := &VocabWordRs{
		ID:      &vocabWord.ID,
		Updated: vocabWord.UpdatedAt.UnixMilli(),
	}

	return c.Status(http.StatusOK).JSON(fext.D(wordRs))
}

func (h *Handler) deleteWord(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	var data RemoveVocabWordRq
	err = c.BodyParser(&data)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgInternal))
	}

	err = h.vocabSvc.DeleteWord(ctx, uid, data.VocabID, data.WordID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgInternal))
	}

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) getWord(c *fiber.Ctx) error {
	ctx := c.Context()

	vid, err := uuid.Parse(c.Query(paramsID))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	wordID, err := uuid.Parse(c.Query(paramsWordID))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	vocabWord, err := h.vocabSvc.GetWord(ctx, vid, wordID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	translates := make([]string, 0, len(vocabWord.Translates))
	for _, translate := range vocabWord.Translates {
		translates = append(translates, translate.Text)
	}

	examples := make([]string, 0, len(vocabWord.Examples))
	for _, example := range vocabWord.Examples {
		examples = append(examples, example.Text)
	}

	wordRs := VocabWordRs{
		ID: &vocabWord.ID,
		Native: &VocabWord{
			ID:            &vocabWord.Native.ID,
			Text:          vocabWord.Native.Text,
			Pronunciation: vocabWord.Native.Pronunciation,
		},
		Definition: vocabWord.Definition,
		Translates: translates,
		Examples:   examples,
	}

	return c.Status(http.StatusOK).JSON(fext.D(wordRs))
}

func (h *Handler) getWords(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, _ := fext.UserIDFromContext(c)

	vid, err := uuid.Parse(c.Query(paramsID))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err))
	}

	vocabWords, err := h.vocabSvc.GetWords(ctx, uid, vid)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	wordsRs := make([]VocabWordRs, 0, len(vocabWords))
	for _, vocabWord := range vocabWords {
		translates := make([]string, 0, len(vocabWord.Translates))
		for _, translate := range vocabWord.Translates {
			translates = append(translates, translate.Text)
		}
		examples := make([]string, 0, len(vocabWord.Examples))
		for _, example := range vocabWord.Examples {
			examples = append(examples, example.Text)
		}

		wordRs := VocabWordRs{
			ID: &vocabWord.ID,
			Native: &VocabWord{
				ID:            &vocabWord.Native.ID,
				Text:          vocabWord.Native.Text,
				Pronunciation: vocabWord.Native.Pronunciation,
			},
			Definition: vocabWord.Definition,
			Translates: translates,
			Examples:   examples,
			Created:    vocabWord.CreatedAt.UnixMilli(),
			Updated:    vocabWord.UpdatedAt.UnixMilli(),
		}

		wordsRs = append(wordsRs, wordRs)
	}

	return c.Status(http.StatusOK).JSON(fext.D(wordsRs))
}

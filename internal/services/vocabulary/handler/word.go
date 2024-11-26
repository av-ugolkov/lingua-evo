package handler

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"
	"unicode/utf8"

	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	entityDict "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	entityExample "github.com/av-ugolkov/lingua-evo/internal/services/example"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	paramsText string = "text"
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

func (h *Handler) addWord(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			fmt.Errorf("word.delivery.Handler.addWord: %v", err)
	}

	var data struct {
		VocabID    uuid.UUID `json:"vocab_id"`
		Native     VocabWord `json:"native"`
		Definition string    `json:"definition,omitempty"`
		Translates []string  `json:"translates,omitempty"`
		Examples   []string  `json:"examples,omitempty"`
	}
	err = c.Bind(&data)
	if err != nil {
		return http.StatusInternalServerError, nil,
			msgerr.New(fmt.Errorf("word.delivery.Handler.addWord: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	if utf8.RuneCountInString(data.Native.Text) > 100 {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("word.delivery.Handler.addWord: word is too long"),
				ErrWordTooLong)
	}

	if utf8.RuneCountInString(data.Native.Pronunciation) > 100 {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("word.delivery.Handler.addWord: pronunciation is too long"),
				ErrPronunciationTooLong)
	}

	if utf8.RuneCountInString(data.Definition) > 100 {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("word.delivery.Handler.addWord: definition is too long"),
				ErrDefinitionTooLong)
	}

	if len(data.Translates) > 10 {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("word.delivery.Handler.addWord: translates more than 10"),
				ErrCountTranslates)
	}

	if len(data.Examples) > 5 {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("word.delivery.Handler.addWord: translates more than 5"),
				ErrCountExamples)
	}

	translateWords := make([]entityDict.DictWord, 0, len(data.Translates))
	for _, translateWord := range data.Translates {
		translateWords = append(translateWords, entityDict.DictWord{
			Text:      translateWord,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
	}

	examples := make([]entityExample.Example, 0, len(data.Examples))
	for _, example := range data.Examples {
		examples = append(examples, entityExample.Example{
			Text:      example,
			CreatedAt: time.Now().UTC(),
		})
	}

	vocabWord, err := h.vocabSvc.AddWord(ctx, userID, entity.VocabWordData{
		VocabID: data.VocabID,
		Native: entityDict.DictWord{
			Text:          data.Native.Text,
			Pronunciation: data.Native.Pronunciation,
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		},
		Definition: data.Definition,
		Translates: translateWords,
		Examples:   examples,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	})
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrDuplicate):
			return http.StatusConflict, nil,
				msgerr.New(fmt.Errorf("word.delivery.Handler.addWord: %v", err),
					ErrWordIsExists)
		default:
			return http.StatusInternalServerError, nil,
				msgerr.New(fmt.Errorf("word.delivery.Handler.addWord: %v", err),
					msgerr.ErrMsgInternal)
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

	return http.StatusCreated, wordRs, nil
}

func (h *Handler) updateWordText(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil, fmt.Errorf("service.vocabulary.Handler.updateWordText: %v", err)
	}

	var data struct {
		ID      uuid.UUID `json:"id"`
		VocabID uuid.UUID `json:"vocab_id"`
		Native  struct {
			ID   uuid.UUID `json:"id"`
			Text string    `json:"text"`
		}
	}
	err = c.Bind(&data)
	if err != nil {
		return http.StatusInternalServerError, nil,
			msgerr.New(fmt.Errorf("service.vocabulary.Handler.updateWordText: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	if utf8.RuneCountInString(data.Native.Text) > 100 {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("service.vocabulary.Handler.updateWordText: word is too long"),
				ErrWordTooLong)
	}

	vocabWord, err := h.vocabSvc.UpdateWordText(ctx, userID, entity.VocabWordData{
		ID:      data.ID,
		VocabID: data.VocabID,
		Native: entityDict.DictWord{
			ID:   data.Native.ID,
			Text: data.Native.Text,
		},
	})
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("service.vocabulary.Handler.updateWordText: %v", err)
	}

	wordRs := &VocabWordRs{
		ID: &vocabWord.ID,
		Native: &VocabWord{
			ID: &vocabWord.NativeID,
		},
		Updated: vocabWord.UpdatedAt.UnixMilli(),
	}

	return http.StatusOK, wordRs, nil
}

func (h *Handler) updateWordPronunciation(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil, fmt.Errorf("service.vocabulary.Handler.updateWordPronunciation: %v", err)
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
	err = c.Bind(&data)
	if err != nil {
		return http.StatusInternalServerError, nil,
			msgerr.New(fmt.Errorf("service.vocabulary.Handler.updateWordPronunciation: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	if utf8.RuneCountInString(data.Native.Pronunciation) > 100 {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("service.vocabulary.Handler.updateWordPronunciation: pronunciation is too long"),
				ErrPronunciationTooLong)
	}

	vocabWord, err := h.vocabSvc.UpdateWordPronunciation(ctx, userID, entity.VocabWordData{
		ID:      data.ID,
		VocabID: data.VocabID,
		Native: entityDict.DictWord{
			ID:            data.Native.ID,
			Text:          data.Native.Text,
			Pronunciation: data.Native.Pronunciation,
		},
	})
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("service.vocabulary.Handler.updateWordPronunciation: %v", err)
	}

	wordRs := &VocabWordRs{
		ID:      &vocabWord.ID,
		Updated: vocabWord.UpdatedAt.UnixMilli(),
	}

	return http.StatusOK, wordRs, nil
}

func (h *Handler) updateWordDefinition(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil, fmt.Errorf("service.vocabulary.Handler.updateWordDefinition: %v", err)
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
	err = c.Bind(&data)
	if err != nil {
		return http.StatusInternalServerError, nil,
			msgerr.New(fmt.Errorf("service.vocabulary.Handler.updateWordDefinition: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	if utf8.RuneCountInString(data.Definition) > 100 {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("service.vocabulary.Handler.updateWordDefinition: definition is too long"),
				ErrDefinitionTooLong)
	}

	vocabWord, err := h.vocabSvc.UpdateWordDefinition(ctx, userID, entity.VocabWordData{
		ID:         data.ID,
		VocabID:    data.VocabID,
		Definition: data.Definition,
	})
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("service.vocabulary.Handler.updateWordDefinition: %v", err)
	}

	wordRs := &VocabWordRs{
		ID:      &vocabWord.ID,
		Updated: vocabWord.UpdatedAt.UnixMilli(),
	}

	return http.StatusOK, wordRs, nil
}

func (h *Handler) updateWordTranslates(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil, fmt.Errorf("service.vocabulary.Handler.updateWordTranslates: %v", err)
	}

	var data struct {
		ID         uuid.UUID `json:"id"`
		VocabID    uuid.UUID `json:"vocab_id"`
		Translates []string  `json:"translates"`
	}
	err = c.Bind(&data)
	if err != nil {
		return http.StatusInternalServerError, nil,
			msgerr.New(fmt.Errorf("service.vocabulary.Handler.updateWordTranslates: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	if len(data.Translates) > 10 {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("service.vocabulary.Handler.updateWordTranslates: translates more than 10"),
				ErrCountTranslates)
	}

	translates := make([]entityDict.DictWord, 0, len(data.Translates))
	for _, tr := range data.Translates {
		if slices.ContainsFunc(translates, func(t entityDict.DictWord) bool {
			return t.Text == tr
		}) {
			return http.StatusBadRequest, nil,
				msgerr.New(fmt.Errorf("service.vocabulary.Handler.updateWordTranslates: duplicate translate"),
					ErrWordIsExists)
		}
		translates = append(translates, entityDict.DictWord{
			Text: tr,
		})
	}

	vocabWord, err := h.vocabSvc.UpdateWordTranslates(ctx, userID, entity.VocabWordData{
		ID:         data.ID,
		VocabID:    data.VocabID,
		Translates: translates,
	})
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("service.vocabulary.Handler.updateWordTranslates: %w", err)
	}

	wordRs := &VocabWordRs{
		ID:      &vocabWord.ID,
		Updated: vocabWord.UpdatedAt.UnixMilli(),
	}

	return http.StatusOK, wordRs, nil
}

func (h *Handler) updateWordExamples(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil, fmt.Errorf("service.vocabulary.Handler.updateWordExamples: %v", err)
	}

	var data struct {
		ID       uuid.UUID `json:"id"`
		VocabID  uuid.UUID `json:"vocab_id"`
		Examples []string  `json:"examples"`
	}
	err = c.Bind(&data)
	if err != nil {
		return http.StatusInternalServerError, nil,
			msgerr.New(fmt.Errorf("service.vocabulary.Handler.updateWordExamples: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	if len(data.Examples) > 5 {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("service.vocabulary.Handler.updateWordExamples: translates more than 5"),
				ErrCountExamples)
	}

	examples := make([]entityExample.Example, 0, len(data.Examples))
	for _, ex := range data.Examples {
		if slices.ContainsFunc(examples, func(e entityExample.Example) bool {
			return e.Text == ex
		}) {
			return http.StatusBadRequest, nil,
				msgerr.New(fmt.Errorf("service.vocabulary.Handler.updateWordExamples: duplicate example"),
					ErrExampleIsExists)
		}
		examples = append(examples, entityExample.Example{
			Text: ex,
		})
	}

	vocabWord, err := h.vocabSvc.UpdateWordExamples(ctx, userID, entity.VocabWordData{
		ID:       data.ID,
		VocabID:  data.VocabID,
		Examples: examples,
	})
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("service.vocabulary.Handler.updateWordExamples: %v", err)
	}

	wordRs := &VocabWordRs{
		ID:      &vocabWord.ID,
		Updated: vocabWord.UpdatedAt.UnixMilli(),
	}

	return http.StatusOK, wordRs, nil
}

func (h *Handler) deleteWord(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			msgerr.New(fmt.Errorf("word.delivery.Handler.deleteWord: %v", err),
				msgerr.ErrMsgInternal)
	}

	var data RemoveVocabWordRq
	err = c.Bind(&data)
	if err != nil {
		return http.StatusInternalServerError, nil,
			msgerr.New(fmt.Errorf("word.delivery.Handler.deleteWord: %v", err),
				msgerr.ErrMsgInternal)
	}

	err = h.vocabSvc.DeleteWord(ctx, uid, data.VocabID, data.WordID)
	if err != nil {
		return http.StatusInternalServerError, nil,
			msgerr.New(fmt.Errorf("word.delivery.Handler.deleteWord: %v", err),
				msgerr.ErrMsgInternal)
	}

	return http.StatusOK, gin.H{}, nil
}

func (h *Handler) getWord(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	vid, err := c.GetQueryUUID(paramsID)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("word.delivery.Handler.getWords: %w", err)
	}

	wordID, err := c.GetQueryUUID(paramsWordID)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("word.delivery.Handler.getWords: %w", err)
	}

	vocabWord, err := h.vocabSvc.GetWord(ctx, vid, wordID)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("word.delivery.Handler.getWords: %w", err)
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

	return http.StatusOK, wordRs, nil
}

func (h *Handler) getWords(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, _ := runtime.UserIDFromContext(ctx)

	vid, err := c.GetQueryUUID(paramsID)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("word.delivery.Handler.getWords: %w", err)
	}

	vocabWords, err := h.vocabSvc.GetWords(ctx, uid, vid)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("word.delivery.Handler.getWords: %w", err)
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

	return http.StatusOK, wordsRs, nil
}

func (h *Handler) getPronunciation(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			fmt.Errorf("word.delivery.Handler.getPronunciation: %w", err)
	}

	text, err := c.GetQuery(paramsText)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("word.delivery.Handler.getPronunciation: %w", err)
	}

	vid, err := c.GetQueryUUID(paramsID)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("word.delivery.Handler.getPronunciation: %w", err)
	}

	pronunciation, err := h.vocabSvc.GetPronunciation(ctx, uid, vid, text)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	vocabWordRs := VocabWordRs{
		Native: &VocabWord{
			Pronunciation: pronunciation,
		},
	}

	return http.StatusOK, vocabWordRs, nil
}

package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	msgerror "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
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
	ErrDescriptionTooLong = "Description length should be less than 100 characters"
	ErrWordIsExists       = "This word is already exists"
)

type (
	VocabWord struct {
		ID            *uuid.UUID `json:"id,omitempty"`
		Text          string     `json:"text,omitempty"`
		Pronunciation string     `json:"pronunciation,omitempty"`
	}

	VocabWordRq struct {
		ID          *uuid.UUID `json:"id,omitempty"`
		VocabID     uuid.UUID  `json:"vocab_id"`
		Native      VocabWord  `json:"native"`
		Description string     `json:"description,omitempty"`
		Translates  []string   `json:"translates,omitempty"`
		Examples    []string   `json:"examples,omitempty"`
	}

	RemoveVocabWordRq struct {
		VocabID uuid.UUID `json:"vocab_id"`
		WordID  uuid.UUID `json:"word_id"`
	}

	VocabWordRs struct {
		ID          *uuid.UUID `json:"id,omitempty"`
		Native      *VocabWord `json:"native,omitempty"`
		Description string     `json:"description,omitempty"`
		Translates  []string   `json:"translates,omitempty"`
		Examples    []string   `json:"examples,omitempty"`
		Created     int64      `json:"created,omitempty"`
		Updated     int64      `json:"updated,omitempty"`
	}
)

func (h *Handler) addWord(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			fmt.Errorf("word.delivery.Handler.addWord: %v", err)
	}

	var data VocabWordRq
	err = c.Bind(&data)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("word.delivery.Handler.addWord: %v", err)
	}

	if utf8.RuneCountInString(data.Description) > 100 {
		return http.StatusBadRequest, nil,
			msgerror.NewError(fmt.Errorf("word.delivery.Handler.addWord - description is too long"),
				ErrDescriptionTooLong)
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
		Description: data.Description,
		Translates:  translateWords,
		Examples:    examples,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrDuplicate):
			return http.StatusConflict, nil,
				msgerror.NewError(fmt.Errorf("word.delivery.Handler.addWord: %v", err), ErrWordIsExists)
		default:
			return http.StatusInternalServerError, nil, fmt.Errorf("word.delivery.Handler.addWord: %v", err)
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

func (h *Handler) updateWord(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil, fmt.Errorf("word.delivery.Handler.updateWord: %v", err)
	}

	var data VocabWordRq
	err = c.Bind(&data)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("word.delivery.Handler.updateWord: %v", err)
	}

	translates := make([]entityDict.DictWord, 0, len(data.Translates))
	for _, tr := range data.Translates {
		translates = append(translates, entityDict.DictWord{
			Text:      tr,
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

	vocabWord, err := h.vocabSvc.UpdateWord(ctx, userID, entity.VocabWordData{
		ID:      *data.ID,
		VocabID: data.VocabID,
		Native: entityDict.DictWord{
			ID:            *data.Native.ID,
			Text:          data.Native.Text,
			Pronunciation: data.Native.Pronunciation,
			UpdatedAt:     time.Now().UTC(),
		},
		Description: data.Description,
		Translates:  translates,
		Examples:    examples,
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("word.delivery.Handler.updateWord: %v", err)
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

func (h *Handler) deleteWord(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	var data RemoveVocabWordRq
	err := c.Bind(&data)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("word.delivery.Handler.deleteWord: %v", err)
	}

	err = h.vocabSvc.DeleteWord(ctx, data.VocabID, data.WordID)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("word.delivery.Handler.deleteWord: %v", err)
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
		Description: vocabWord.Description,
		Translates:  translates,
		Examples:    examples,
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
			Description: vocabWord.Description,
			Translates:  translates,
			Examples:    examples,
			Created:     vocabWord.CreatedAt.UnixMilli(),
			Updated:     vocabWord.UpdatedAt.UnixMilli(),
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

package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	vocabulary "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/service"
	"github.com/av-ugolkov/lingua-evo/runtime"
	"github.com/av-ugolkov/lingua-evo/runtime/access"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	paramsID            string = "id"
	paramsWordID        string = "word_id"
	paramsVocabName     string = "name"
	paramsPage          string = "page"
	paramsPerPage       string = "per_page"
	paramsSearch        string = "search"
	paramsSort          string = "sort"
	paramsOrder         string = "order"
	paramsNativeLang    string = "native_lang"
	paramsTranslateLang string = "translate_lang"
	paramsUserID        string = "user_id"
	paramsLimitWords    string = "limit_words"
)

type (
	VocabularyRq struct {
		ID            uuid.UUID `json:"id,omitempty"`
		Name          string    `json:"name,omitempty"`
		Access        uint8     `json:"access_id,omitempty"`
		NativeLang    string    `json:"native_lang,omitempty"`
		TranslateLang string    `json:"translate_lang,omitempty"`
		Description   string    `json:"description,omitempty"`
	}

	VocabularyRs struct {
		ID            *uuid.UUID `json:"id,omitempty"`
		UserID        *uuid.UUID `json:"user_id,omitempty"`
		Name          string     `json:"name,omitempty"`
		AccessID      *uint8     `json:"access_id,omitempty"`
		NativeLang    string     `json:"native_lang,omitempty"`
		TranslateLang string     `json:"translate_lang,omitempty"`
		Description   string     `json:"description,omitempty"`
		UserName      string     `json:"user_name,omitempty"`
		Editable      *bool      `json:"editable,omitempty"`
		WordsCount    *uint      `json:"words_count,omitempty"`
		CreatedAt     *time.Time `json:"created_at,omitempty"`
		UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	}

	VocabByUserRs struct {
		ID            uuid.UUID `json:"id,omitempty"`
		UserID        uuid.UUID `json:"user_id,omitempty"`
		Name          string    `json:"name,omitempty"`
		AccessID      *uint8    `json:"access_id,omitempty"`
		NativeLang    string    `json:"native_lang,omitempty"`
		TranslateLang string    `json:"translate_lang,omitempty"`
		Description   string    `json:"description,omitempty"`
		UserName      string    `json:"user_name,omitempty"`
		WordsCount    *uint     `json:"words_count,omitempty"`
		Editable      bool      `json:"editable,omitempty"`
		Notification  bool      `json:"notification,omitempty"`
	}

	VocabularyWithWords struct {
		VocabularyRs
		Words []string `json:"words,omitempty"`
	}
)

type Handler struct {
	vocabSvc *vocabulary.Service
}

func Create(r *fiber.App, vocabSvc *vocabulary.Service) {
	h := newHandler(vocabSvc)

	r.Post(handler.Vocabulary, middleware.Auth(h.userAddVocabulary))
	r.Delete(handler.Vocabulary, middleware.Auth(h.userDeleteVocabulary))
	r.Put(handler.Vocabulary, middleware.Auth(h.userEditVocabulary))
	r.Get(handler.UserVocabularies, middleware.Auth(h.userGetVocabularies))
	r.Get(handler.Vocabularies, middleware.OptionalAuth(h.getVocabularies))
	r.Get(handler.VocabulariesByUser, middleware.OptionalAuth(h.getVocabulariesByUser))
	r.Get(handler.VocabularyInfo, middleware.OptionalAuth(h.getVocabularyInfo))
	r.Post(handler.VocabularyCopy, middleware.Auth(h.copyVocabulary))
	r.Get(handler.VocabulariesRecommended, middleware.OptionalAuth(h.getRecommendedVocabularies))

	r.Get(handler.VocabularyWord, middleware.Auth(h.getWord))
	r.Post(handler.VocabularyWord, middleware.Auth(h.addWord))
	r.Delete(handler.VocabularyWord, middleware.Auth(h.deleteWord))
	r.Post(handler.VocabularyWordText, middleware.Auth(h.updateWordText))
	r.Post(handler.VocabularyWordPronunciation, middleware.Auth(h.updateWordPronunciation))
	r.Post(handler.VocabularyWordDefinition, middleware.Auth(h.updateWordDefinition))
	r.Post(handler.VocabularyWordTranslates, middleware.Auth(h.updateWordTranslates))
	r.Post(handler.VocabularyWordExamples, middleware.Auth(h.updateWordExamples))
	r.Get(handler.VocabularyWords, middleware.OptionalAuth(h.getWords))

	r.Get(handler.VocabularyAccessForUser, middleware.Auth(h.getAccessForUser))
	r.Post(handler.VocabularyAccessForUser, middleware.Auth(h.addAccessForUser))
	r.Delete(handler.VocabularyAccessForUser, middleware.Auth(h.removeAccessForUser))
	r.Patch(handler.VocabularyAccessForUser, middleware.Auth(h.updateAccessForUser))
}

func newHandler(vocabSvc *vocabulary.Service) *Handler {
	return &Handler{
		vocabSvc: vocabSvc,
	}
}

func (h *Handler) getVocabularies(c *fiber.Ctx) error {
	ctx := c.Context()
	uid, _ := runtime.UserIDFromContext(ctx)

	page := c.QueryInt(paramsPage)
	if page == 0 {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("vocabulary.Handler.getVocabularies: not found query [%s]", paramsPage))
	}

	itemsPerPage := c.QueryInt(paramsPerPage)
	if itemsPerPage == 0 {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("vocabulary.Handler.getVocabularies: not found query [%s]", paramsPerPage))
	}

	typeSort := c.QueryInt(paramsSort, -1)
	if typeSort == -1 {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("vocabulary.Handler.getVocabularies: not found query [%s]", paramsSort))
	}

	order := c.QueryInt(paramsOrder, -1)
	if order == -1 {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("vocabulary.Handler.getVocabularies: not found query [%s]", paramsOrder))
	}

	search := c.Query(paramsSearch)
	nativeLang := c.Query(paramsNativeLang)
	if nativeLang == runtime.EmptyString {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("vocabulary.Handler.getVocabularies: not found query [%s]", paramsNativeLang))
	}

	translateLang := c.Query(paramsTranslateLang)
	if translateLang == runtime.EmptyString {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("vocabulary.Handler.getVocabularies: not found query [%s]", paramsNativeLang))
	}

	limitWords := c.QueryInt(paramsLimitWords)
	if limitWords == 0 {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("vocabulary.Handler.getVocabularies: not found query [%s]", paramsLimitWords))
	}

	vocabularies, totalCount, err := h.vocabSvc.GetVocabularies(ctx, uid, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang, limitWords)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.getVocabularies: %v", err))
	}

	vocabsWithWordsRs := make([]VocabularyWithWords, 0, len(vocabularies))
	for _, vocab := range vocabularies {
		vocabsWithWordsRs = append(vocabsWithWordsRs, VocabularyWithWords{
			VocabularyRs: VocabularyRs{
				ID:            &vocab.ID,
				UserID:        &vocab.UserID,
				UserName:      vocab.UserName,
				Name:          vocab.Name,
				AccessID:      &vocab.Access,
				NativeLang:    vocab.NativeLang,
				TranslateLang: vocab.TranslateLang,
				Description:   vocab.Description,
				WordsCount:    &vocab.WordsCount,
				CreatedAt:     &vocab.CreatedAt,
				UpdatedAt:     &vocab.UpdatedAt,
			},
			Words: vocab.Words,
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"vocabularies": vocabsWithWordsRs,
		"total_count":  totalCount,
	})
}

func (h *Handler) getVocabulariesByUser(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, _ := runtime.UserIDFromContext(ctx)

	owner, err := uuid.Parse(c.Query(paramsUserID))
	if err == nil {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("vocabulary.Handler.getVocabulariesByUser: %v", err))
	}

	vocabs, err := h.vocabSvc.GetVocabulariesByUser(ctx, uid, owner, []access.Type{access.Public, access.Subscribers})
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.getVocabulariesByUser: %v", err))
	}

	vocabulariesRs := make([]VocabByUserRs, 0, len(vocabs))
	for _, vocab := range vocabs {
		vocabulariesRs = append(vocabulariesRs, VocabByUserRs{
			ID:            vocab.ID,
			Name:          vocab.Name,
			UserID:        vocab.UserID,
			UserName:      vocab.UserName,
			AccessID:      &vocab.Access,
			NativeLang:    vocab.NativeLang,
			TranslateLang: vocab.TranslateLang,
			Description:   vocab.Description,
			WordsCount:    &vocab.WordsCount,
			Editable:      vocab.Editable,
			Notification:  vocab.Notification,
		})
	}

	return c.Status(http.StatusOK).JSON(vocabulariesRs)
}

func (h *Handler) getVocabularyInfo(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, _ := runtime.UserIDFromContext(ctx)

	vid, err := uuid.Parse(c.Query(paramsID))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("vocabulary.Handler.getVocabulary: %v", err))
	}

	vocab, err := h.vocabSvc.GetVocabularyInfo(ctx, uid, vid)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.getVocabulary: %v", err))
	}
	if vocab.ID == uuid.Nil {
		return fiber.NewError(http.StatusNotFound,
			fmt.Sprintf("vocabulary.Handler.getVocabulary - vocabulary not found: %v", err))
	}

	vocabRs := VocabularyRs{
		ID:            &vocab.ID,
		AccessID:      &vocab.Access,
		UserID:        &vocab.UserID,
		UserName:      vocab.UserName,
		Name:          vocab.Name,
		NativeLang:    vocab.NativeLang,
		TranslateLang: vocab.TranslateLang,
		Description:   vocab.Description,
		Editable:      &vocab.Editable,
		WordsCount:    &vocab.WordsCount,
		CreatedAt:     &vocab.CreatedAt,
		UpdatedAt:     &vocab.UpdatedAt,
	}

	return c.Status(http.StatusOK).JSON(vocabRs)
}

func (h *Handler) copyVocabulary(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized,
			fmt.Sprintf("vocabulary.Handler.copyVocabulary: %v", err))
	}

	vid, err := uuid.Parse(c.Query(paramsID))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("vocabulary.Handler.copyVocabulary: %v", err))
	}

	err = h.vocabSvc.CopyVocab(ctx, uid, vid)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.copyVocabulary: %v", err))
	}

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) getRecommendedVocabularies(c *fiber.Ctx) error {
	ctx := c.Context()
	uid, _ := runtime.UserIDFromContext(ctx)

	vocabs, err := h.vocabSvc.GetRecommendedVocabularies(ctx, uid)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.getRecommendedVocabularies: %v", err))
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
			Description:   vocab.Description,
			WordsCount:    &vocab.WordsCount,
		})
	}

	return c.Status(http.StatusOK).JSON(vocabulariesRs)
}

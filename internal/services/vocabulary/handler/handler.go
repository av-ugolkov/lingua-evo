package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	vocabulary "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/service"
	"github.com/av-ugolkov/lingua-evo/runtime"
	"github.com/av-ugolkov/lingua-evo/runtime/access"

	"github.com/gin-gonic/gin"
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
		Tags          []string  `json:"tags,omitempty"`
	}

	VocabularyRs struct {
		ID            uuid.UUID `json:"id,omitempty"`
		UserID        uuid.UUID `json:"user_id,omitempty"`
		Name          string    `json:"name,omitempty"`
		AccessID      *uint8    `json:"access_id,omitempty"`
		NativeLang    string    `json:"native_lang,omitempty"`
		TranslateLang string    `json:"translate_lang,omitempty"`
		Description   string    `json:"description,omitempty"`
		Tags          []string  `json:"tags,omitempty"`
		UserName      string    `json:"user_name,omitempty"`
		Editable      *bool     `json:"editable,omitempty"`
		WordsCount    *uint     `json:"words_count,omitempty"`
		CreatedAt     time.Time `json:"created_at,omitempty"`
		UpdatedAt     time.Time `json:"updated_at,omitempty"`
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

func Create(r *ginext.Engine, vocabSvc *vocabulary.Service) {
	h := newHandler(vocabSvc)

	r.POST(handler.Vocabulary, middleware.Auth(h.userAddVocabulary))
	r.DELETE(handler.Vocabulary, middleware.Auth(h.userDeleteVocabulary))
	r.PUT(handler.Vocabulary, middleware.Auth(h.userEditVocabulary))
	r.GET(handler.UserVocabularies, middleware.Auth(h.userGetVocabularies))
	r.GET(handler.Vocabularies, middleware.OptionalAuth(h.getVocabularies))
	r.GET(handler.VocabulariesByUser, middleware.OptionalAuth(h.getVocabulariesByUser))
	r.GET(handler.VocabularyInfo, middleware.OptionalAuth(h.getVocabularyInfo))
	r.POST(handler.VocabularyCopy, middleware.Auth(h.copyVocabulary))
	r.GET(handler.VocabulariesRecommended, middleware.OptionalAuth(h.getRecommendedVocabularies))

	r.GET(handler.VocabularyWord, middleware.Auth(h.getWord))
	r.POST(handler.VocabularyWord, middleware.Auth(h.addWord))
	r.DELETE(handler.VocabularyWord, middleware.Auth(h.deleteWord))
	r.POST(handler.VocabularyWordUpdateFull, middleware.Auth(h.updateWord))
	r.POST(handler.VocabularyWordUpdateText, middleware.Auth(h.updateWordText))
	r.POST(handler.VocabularyWordUpdatePronunciation, middleware.Auth(h.updateWordPronunciation))
	r.POST(handler.VocabularyWordUpdateDefinition, middleware.Auth(h.updateWordDefinition))
	r.POST(handler.VocabularyWordUpdateTranslates, middleware.Auth(h.updateWordTranslates))
	r.POST(handler.VocabularyWordUpdateExamples, middleware.Auth(h.updateWordExamples))
	r.GET(handler.VocabularyWords, middleware.OptionalAuth(h.getWords))
	r.GET(handler.WordPronunciation, middleware.Auth(h.getPronunciation))

	r.GET(handler.VocabularyAccessForUser, middleware.Auth(h.getAccessForUser))
	r.POST(handler.VocabularyAccessForUser, middleware.Auth(h.addAccessForUser))
	r.DELETE(handler.VocabularyAccessForUser, middleware.Auth(h.removeAccessForUser))
	r.PATCH(handler.VocabularyAccessForUser, middleware.Auth(h.updateAccessForUser))
}

func newHandler(vocabSvc *vocabulary.Service) *Handler {
	return &Handler{
		vocabSvc: vocabSvc,
	}
}

func (h *Handler) getVocabularies(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()
	uid, _ := runtime.UserIDFromContext(ctx)

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

	limitWords, err := c.GetQueryInt(paramsLimitWords)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies: %v", err)
	}

	vocabularies, totalCount, err := h.vocabSvc.GetVocabularies(ctx, uid, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang, limitWords)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies: %v", err)
	}

	vocabsWithWordsRs := make([]VocabularyWithWords, 0, len(vocabularies))
	for _, vocab := range vocabularies {
		tags := make([]string, 0, len(vocab.Tags))
		for _, tag := range vocab.Tags {
			tags = append(tags, tag.Text)
		}

		vocabsWithWordsRs = append(vocabsWithWordsRs, VocabularyWithWords{
			VocabularyRs: VocabularyRs{
				ID:            vocab.ID,
				UserID:        vocab.UserID,
				UserName:      vocab.UserName,
				Name:          vocab.Name,
				AccessID:      &vocab.Access,
				NativeLang:    vocab.NativeLang,
				TranslateLang: vocab.TranslateLang,
				Description:   vocab.Description,
				WordsCount:    &vocab.WordsCount,
				Tags:          tags,
				CreatedAt:     vocab.CreatedAt,
				UpdatedAt:     vocab.UpdatedAt,
			},
			Words: vocab.Words,
		})
	}

	return http.StatusOK, gin.H{
		"vocabularies": vocabsWithWordsRs,
		"total_count":  totalCount,
	}, nil
}

func (h *Handler) getVocabulariesByUser(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, _ := runtime.UserIDFromContext(ctx)

	owner, err := c.GetQueryUUID(paramsUserID)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabulariesByUser: %v", err)
	}

	vocabs, err := h.vocabSvc.GetVocabulariesByUser(ctx, uid, owner, []access.Type{access.Public, access.Subscribers})
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabulariesByUser: %v", err)
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

	return http.StatusOK, vocabulariesRs, nil
}

func (h *Handler) getVocabularyInfo(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, _ := runtime.UserIDFromContext(ctx)

	vid, err := c.GetQueryUUID(paramsID)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabulary: %v", err)
	}

	vocab, err := h.vocabSvc.GetVocabularyInfo(ctx, uid, vid)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("vocabulary.delivery.Handler.getVocabulary: %v", err)
	}
	if vocab.ID == uuid.Nil {
		return http.StatusNotFound, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabulary - vocabulary not found: %v", err)
	}

	tags := make([]string, 0, len(vocab.Tags))
	for _, tag := range vocab.Tags {
		tags = append(tags, tag.Text)
	}

	vocabRs := VocabularyRs{
		ID:            vocab.ID,
		AccessID:      &vocab.Access,
		UserID:        vocab.UserID,
		UserName:      vocab.UserName,
		Name:          vocab.Name,
		NativeLang:    vocab.NativeLang,
		TranslateLang: vocab.TranslateLang,
		Description:   vocab.Description,
		Editable:      &vocab.Editable,
		Tags:          tags,
		WordsCount:    &vocab.WordsCount,
		CreatedAt:     vocab.CreatedAt,
		UpdatedAt:     vocab.UpdatedAt,
	}

	return http.StatusOK, vocabRs, nil
}

func (h *Handler) copyVocabulary(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			fmt.Errorf("vocabulary.delivery.Handler.copyVocabulary: %v", err)
	}

	vid, err := c.GetQueryUUID(paramsID)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("vocabulary.delivery.Handler.copyVocabulary: %v", err)
	}

	err = h.vocabSvc.CopyVocab(ctx, uid, vid)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.copyVocabulary: %v", err)
	}

	return http.StatusOK, gin.H{}, nil
}

func (h *Handler) getRecommendedVocabularies(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()
	uid, _ := runtime.UserIDFromContext(ctx)

	vocabs, err := h.vocabSvc.GetRecommendedVocabularies(ctx, uid)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getRecommendedVocabularies: %v", err)
	}

	vocabulariesRs := make([]VocabularyRs, 0, len(vocabs))
	for _, vocab := range vocabs {
		tags := make([]string, 0, len(vocab.Tags))
		for _, tag := range vocab.Tags {
			tags = append(tags, tag.Text)
		}

		vocabulariesRs = append(vocabulariesRs, VocabularyRs{
			ID:            vocab.ID,
			UserID:        vocab.UserID,
			Name:          vocab.Name,
			AccessID:      &vocab.Access,
			NativeLang:    vocab.NativeLang,
			TranslateLang: vocab.TranslateLang,
			Description:   vocab.Description,
			Tags:          tags,
			WordsCount:    &vocab.WordsCount,
		})
	}

	return http.StatusOK, vocabulariesRs, nil
}

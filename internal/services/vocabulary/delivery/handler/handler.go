package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/delivery/handler/gin"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	vocabulary "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/service"
	"github.com/av-ugolkov/lingua-evo/runtime"
	"github.com/av-ugolkov/lingua-evo/runtime/access"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	paramsID            string = "id"
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

func Create(r *gin.Engine, vocabSvc *vocabulary.Service) {
	h := newHandler(vocabSvc)
	h.register(r)
}

func newHandler(vocabSvc *vocabulary.Service) *Handler {
	return &Handler{
		vocabSvc: vocabSvc,
	}
}

func (h *Handler) register(r *gin.Engine) {
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
	r.POST(handler.VocabularyWordUpdate, middleware.Auth(h.updateWord))
	r.GET(handler.VocabularyWords, middleware.OptionalAuth(h.getWords))
	r.GET(handler.WordPronunciation, middleware.Auth(h.getPronunciation))

	r.GET(handler.VocabularyAccessForUser, middleware.Auth(h.getAccessForUser))
	r.POST(handler.VocabularyAccessForUser, middleware.Auth(h.addAccessForUser))
	r.DELETE(handler.VocabularyAccessForUser, middleware.Auth(h.removeAccessForUser))
	r.PATCH(handler.VocabularyAccessForUser, middleware.Auth(h.updateAccessForUser))
}

func (h *Handler) getVocabularies(c *gin.Context) {
	ctx := c.Request.Context()
	userID, _ := runtime.UserIDFromContext(ctx)

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

	limitWords, err := ginExt.GetQueryInt(c, paramsLimitWords)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies - get query [limit_words]: %v", err))
		return
	}

	vocabularies, totalCount, err := h.vocabSvc.GetVocabularies(ctx, userID, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang, limitWords)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies: %v", err))
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

	c.JSON(http.StatusOK, gin.H{
		"vocabularies": vocabsWithWordsRs,
		"total_count":  totalCount,
	})
}

func (h *Handler) getVocabulariesByUser(c *gin.Context) {
	ctx := c.Request.Context()

	uid, _ := runtime.UserIDFromContext(ctx)

	owner, err := ginExt.GetQueryUUID(c, paramsUserID)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabulariesByUser - get query [user_id]: %v", err))
		return
	}

	vocabs, err := h.vocabSvc.GetVocabulariesByUser(ctx, uid, owner, []access.Type{access.Public, access.Subscribers})
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabulariesByUser: %v", err))
		return
	}

	vocabulariesRs := make([]VocabByUserRs, 0, len(vocabs))
	for _, vocab := range vocabs {
		tags := make([]string, 0, len(vocab.Tags))
		for _, tag := range vocab.Tags {
			tags = append(tags, tag.Text)
		}

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

	c.JSON(http.StatusOK, vocabulariesRs)
}

func (h *Handler) getVocabularyInfo(c *gin.Context) {
	ctx := c.Request.Context()

	uid, _ := runtime.UserIDFromContext(ctx)

	vid, err := ginExt.GetQueryUUID(c, paramsID)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabulary - get query [name]: %v", err))
		return
	}

	vocab, err := h.vocabSvc.GetVocabularyInfo(ctx, uid, vid)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError, err)
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

	c.JSON(http.StatusOK, vocabRs)
}

func (h *Handler) copyVocabulary(c *gin.Context) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ginExt.SendError(c, http.StatusUnauthorized,
			fmt.Errorf("vocabulary.delivery.Handler.copyVocabulary - get user id: %v", err))
		return
	}

	vid, err := ginExt.GetQueryUUID(c, paramsID)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("vocabulary.delivery.Handler.copyVocabulary - get query [id]: %v", err))
		return
	}

	err = h.vocabSvc.CopyVocab(ctx, uid, vid)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.copyVocabulary: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) getRecommendedVocabularies(c *gin.Context) {
	ctx := c.Request.Context()
	uid, _ := runtime.UserIDFromContext(ctx)

	vocabs, err := h.vocabSvc.GetRecommendedVocabularies(ctx, uid)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getRecommendedVocabularies: %v", err))
		return
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

	c.JSON(http.StatusOK, vocabulariesRs)
}

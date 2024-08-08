package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/delivery/handler/gin"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	vocabulary "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/service"
	"github.com/av-ugolkov/lingua-evo/runtime"
)

const (
	paramsVocabName     string = "name"
	paramsVocabID       string = "id"
	paramsPage          string = "page"
	paramsPerPage       string = "per_page"
	paramsSearch        string = "search"
	paramsSort          string = "sort"
	paramsOrder         string = "order"
	paramsNativeLang    string = "native_lang"
	paramsTranslateLang string = "translate_lang"
	paramsUserID        string = "user_id"
)

type (
	VocabularyRq struct {
		Name          string   `json:"name"`
		Access        uint8    `json:"access_id"`
		NativeLang    string   `json:"native_lang"`
		TranslateLang string   `json:"translate_lang"`
		Description   string   `json:"description"`
		Tags          []string `json:"tags"`
	}

	VocabularyIDRs struct {
		ID uuid.UUID `json:"id"`
	}

	VocabularyRs struct {
		ID            uuid.UUID `json:"id"`
		UserID        uuid.UUID `json:"user_id"`
		Name          string    `json:"name"`
		AccessID      uint8     `json:"access_id"`
		NativeLang    string    `json:"native_lang"`
		TranslateLang string    `json:"translate_lang"`
		Description   string    `json:"description"`
		Tags          []string  `json:"tags"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}

	VocabularyWithUserRs struct {
		ID            uuid.UUID `json:"id"`
		UserID        uuid.UUID `json:"user_id"`
		UserName      string    `json:"user_name"`
		Name          string    `json:"name"`
		AccessID      uint8     `json:"access_id"`
		NativeLang    string    `json:"native_lang"`
		TranslateLang string    `json:"translate_lang"`
		Description   string    `json:"description"`
		WordsCount    uint      `json:"words_count"`
		Tags          []string  `json:"tags"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}

	VocabularyEditRq struct {
		ID     uuid.UUID `json:"id"`
		Name   string    `json:"name"`
		Access uint8     `json:"access_id"`
	}
)

type Handler struct {
	vocabularySvc *vocabulary.Service
}

func Create(r *gin.Engine, vocabularySvc *vocabulary.Service) {
	h := newHandler(vocabularySvc)
	h.register(r)
}

func newHandler(vocabularySvc *vocabulary.Service) *Handler {
	return &Handler{
		vocabularySvc: vocabularySvc,
	}
}

func (h *Handler) register(r *gin.Engine) {
	r.POST(handler.UserVocabulary, middleware.Auth(h.userAddVocabulary))
	r.DELETE(handler.UserVocabulary, middleware.Auth(h.userDeleteVocabulary))
	r.GET(handler.Vocabulary, middleware.OptionalAuth(h.getVocabulary))
	r.PUT(handler.UserVocabulary, middleware.Auth(h.userEditVocabulary))
	r.GET(handler.UserVocabularies, middleware.Auth(h.userGetVocabularies))
	r.GET(handler.Vocabularies, middleware.OptionalAuth(h.getVocabularies))
	r.GET(handler.VocabulariesByUser, middleware.OptionalAuth(h.getVocabulariesByUser))
	r.GET(handler.VocabularyInfo, middleware.OptionalAuth(h.getVocabularyInfo))
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

	vocabularies, totalCount, err := h.vocabularySvc.GetVocabularies(ctx, userID, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabularies: %v", err))
	}

	vocabulariesRs := make([]VocabularyWithUserRs, 0, len(vocabularies))
	for _, vocab := range vocabularies {
		tags := make([]string, 0, len(vocab.Tags))
		for _, tag := range vocab.Tags {
			tags = append(tags, tag.Text)
		}

		vocabulariesRs = append(vocabulariesRs, VocabularyWithUserRs{
			ID:            vocab.ID,
			UserID:        vocab.UserID,
			UserName:      vocab.UserName,
			Name:          vocab.Name,
			AccessID:      vocab.Access,
			NativeLang:    vocab.NativeLang,
			TranslateLang: vocab.TranslateLang,
			Description:   vocab.Description,
			WordsCount:    vocab.WordsCount,
			Tags:          tags,
			CreatedAt:     vocab.CreatedAt,
			UpdatedAt:     vocab.UpdatedAt,
		})
	}

	var rs struct {
		Vocabularies []VocabularyWithUserRs `json:"vocabularies"`
		TotalCount   int                    `json:"total_count"`
	}
	rs.Vocabularies = vocabulariesRs
	rs.TotalCount = totalCount

	c.JSON(http.StatusOK, rs)
}

func (h *Handler) getVocabulary(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _ := runtime.UserIDFromContext(ctx)

	vocabID, err := ginExt.GetQueryUUID(c, paramsVocabID)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabulary - get query [name]: %v", err))
		return
	}

	vocab, err := h.vocabularySvc.GetVocabulary(ctx, userID, vocabID)
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
		AccessID:      vocab.Access,
		UserID:        vocab.UserID,
		Name:          vocab.Name,
		NativeLang:    vocab.NativeLang,
		TranslateLang: vocab.TranslateLang,
		Description:   vocab.Description,
		Tags:          tags,
		CreatedAt:     vocab.CreatedAt,
		UpdatedAt:     vocab.UpdatedAt,
	}

	c.JSON(http.StatusOK, vocabRs)
}

func (h *Handler) getVocabulariesByUser(c *gin.Context) {
	ctx := c.Request.Context()

	uid, err := ginExt.GetQueryUUID(c, paramsUserID)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabulariesByUser - get query [user_id]: %v", err))
		return
	}

	vocabs, err := h.vocabularySvc.GetVocabulariesByUser(ctx, uid, []entity.AccessVocab{entity.AccessPublic, entity.AccessSubscribers})
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabulariesByUser: %v", err))
		return
	}

	vocabulariesRs := make([]VocabularyWithUserRs, 0, len(vocabs))
	for _, vocab := range vocabs {
		tags := make([]string, 0, len(vocab.Tags))
		for _, tag := range vocab.Tags {
			tags = append(tags, tag.Text)
		}

		vocabulariesRs = append(vocabulariesRs, VocabularyWithUserRs{
			ID:            vocab.ID,
			UserID:        vocab.UserID,
			Name:          vocab.Name,
			AccessID:      vocab.Access,
			NativeLang:    vocab.NativeLang,
			TranslateLang: vocab.TranslateLang,
			Description:   vocab.Description,
			WordsCount:    vocab.WordsCount,
			Tags:          tags,
			CreatedAt:     vocab.CreatedAt,
			UpdatedAt:     vocab.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, vocabulariesRs)
}

func (h *Handler) getVocabularyInfo(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _ := runtime.UserIDFromContext(ctx)

	vocabID, err := ginExt.GetQueryUUID(c, paramsVocabID)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getVocabulary - get query [name]: %v", err))
		return
	}

	vocab, err := h.vocabularySvc.GetVocabularyInfo(ctx, userID, vocabID)
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

	vocabRs := VocabularyWithUserRs{
		ID:            vocab.ID,
		AccessID:      vocab.Access,
		UserID:        vocab.UserID,
		UserName:      vocab.UserName,
		Name:          vocab.Name,
		NativeLang:    vocab.NativeLang,
		TranslateLang: vocab.TranslateLang,
		Description:   vocab.Description,
		Tags:          tags,
		WordsCount:    vocab.WordsCount,
		CreatedAt:     vocab.CreatedAt,
		UpdatedAt:     vocab.UpdatedAt,
	}

	c.JSON(http.StatusOK, vocabRs)
}

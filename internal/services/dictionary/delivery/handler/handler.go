package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	dictionarySvc "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
)

const (
	QueryParamText     = "text"
	QueryParamLangCode = "lang_code"
)

type (
	WordRq struct {
		ID            *uuid.UUID `json:"id,omitempty"`
		Text          string     `json:"text,omitempty"`
		Pronunciation string     `json:"pronunciation,omitempty"`
		LangCode      string     `json:"lang_code,omitempty"`
		CreatedAt     *time.Time `json:"created_at,omitempty"`
		UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	}

	WordRs struct {
		ID            *uuid.UUID `json:"id,omitempty"`
		Text          string     `json:"text,omitempty"`
		Pronunciation string     `json:"pronunciation,omitempty"`
		LangCode      string     `json:"lang_code,omitempty"`
		Creator       *uuid.UUID `json:"creator,omitempty"`
		Moderator     *uuid.UUID `json:"moderator,omitempty"`
		CreatedAt     *time.Time `json:"created_at,omitempty"`
		UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	}
)

type Handler struct {
	dictSvc *dictionarySvc.Service
}

func Create(r *ginext.Engine, dictSvc *dictionarySvc.Service) {
	h := newHandler(dictSvc)

	r.POST(handler.DictionaryWord, middleware.Auth(h.addWord))
	r.GET(handler.DictionaryWord, h.getWord)
	r.GET(handler.GetRandomWord, h.getRandomWord)
}

func newHandler(dictSvc *dictionarySvc.Service) *Handler {
	return &Handler{
		dictSvc: dictSvc,
	}
}

func (h *Handler) addWord(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	var data WordRq
	if err := c.Bind(&data); err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("dictionary.delivery.Handler.addWord - check body: %v", err)
	}

	words, err := h.dictSvc.GetOrAddWords(ctx, []entity.DictWord{
		{
			Text:          data.Text,
			Pronunciation: data.Pronunciation,
			LangCode:      data.LangCode,
			UpdatedAt:     time.Now().UTC(),
			CreatedAt:     time.Now().UTC(),
		},
	})
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("dictionary.delivery.Handler.addWord: %v", err)
	}

	wordRs := &WordRs{
		ID:            &words[0].ID,
		Text:          words[0].Text,
		Pronunciation: words[0].Pronunciation,
		CreatedAt:     &words[0].CreatedAt,
		UpdatedAt:     &words[0].UpdatedAt,
	}

	return http.StatusOK, wordRs, nil
}

func (h *Handler) getWord(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	text, err := c.GetQuery(QueryParamText)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("dictionary.delivery.Handler.getWord: %w", err)
	}

	langCode, err := c.GetQuery(QueryParamLangCode)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("dictionary.delivery.Handler.getWord: %w", err)
	}

	if langCode == runtime.EmptyString {
		return http.StatusBadRequest, nil,
			fmt.Errorf("dictionary.delivery.Handler.getWord - empty lang code")
	}

	wordIDs, err := h.dictSvc.GetWordsByText(ctx, []entity.DictWord{
		{
			Text:     text,
			LangCode: langCode,
		},
	})
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("dictionary.delivery.Handler.getWord: %v", err)
	}

	return http.StatusOK, &WordRs{
		ID: &wordIDs[0].ID,
	}, nil
}

func (h *Handler) getRandomWord(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	langCode, err := c.GetQuery(QueryParamLangCode)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("dictionary.delivery.Handler.getRandomWord: %w", err)
	}

	word, err := h.dictSvc.GetRandomWord(ctx, langCode)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("dictionary.delivery.Handler.getRandomWord: %v", err)
	}

	return http.StatusOK, &WordRs{
		Text:          word.Text,
		LangCode:      word.LangCode,
		Pronunciation: word.Pronunciation,
	}, nil
}

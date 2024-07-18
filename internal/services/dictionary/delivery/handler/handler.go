package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/delivery/handler/gin"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	dictionarySvc "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gin-gonic/gin"
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

func Create(r *gin.Engine, dictSvc *dictionarySvc.Service) {
	h := newHandler(dictSvc)
	h.register(r)
}

func newHandler(dictSvc *dictionarySvc.Service) *Handler {
	return &Handler{
		dictSvc: dictSvc,
	}
}

func (h *Handler) register(r *gin.Engine) {
	r.POST(handler.DictionaryWord, middleware.Auth(h.addWord))
	r.GET(handler.DictionaryWord, h.getWord)
	r.GET(handler.GetRandomWord, h.getRandomWord)
}

func (h *Handler) addWord(c *gin.Context) {
	ctx := c.Request.Context()
	var data WordRq

	if err := c.Bind(&data); err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("dictionary.delivery.Handler.addWord - check body: %v", err))
		return
	}

	words, err := h.dictSvc.GetOrAddWords(ctx, []entity.DictWord{
		{
			ID:            uuid.New(),
			Text:          data.Text,
			Pronunciation: data.Pronunciation,
			LangCode:      data.LangCode,
			UpdatedAt:     time.Now().UTC(),
			CreatedAt:     time.Now().UTC(),
		},
	})
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("dictionary.delivery.Handler.addWord: %v", err))
		return
	}

	wordRs := &WordRs{
		ID:            &words[0].ID,
		Text:          words[0].Text,
		Pronunciation: words[0].Pronunciation,
		CreatedAt:     &words[0].CreatedAt,
		UpdatedAt:     &words[0].UpdatedAt,
	}

	c.JSON(http.StatusOK, wordRs)
}

func (h *Handler) getWord(c *gin.Context) {
	ctx := c.Request.Context()

	text, ok := c.GetQuery(QueryParamText)
	if !ok {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("dictionary.delivery.Handler.getWord - not found query param [text]"))
		return
	}

	langCode, ok := c.GetQuery(QueryParamLangCode)
	if !ok {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("dictionary.delivery.Handler.getWord - not found query param [lang_code]"))
		return
	}

	if langCode == runtime.EmptyString {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("dictionary.delivery.Handler.getWord - empty lang code"))
		return
	}

	wordIDs, err := h.dictSvc.GetWordsByText(ctx, []entity.DictWord{
		{
			Text:     text,
			LangCode: langCode,
		},
	})
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("dictionary.delivery.Handler.getWord: %v", err))
		return
	}

	wordRs := &WordRs{
		ID: &wordIDs[0].ID,
	}

	c.JSON(http.StatusOK, wordRs)
}

func (h *Handler) getRandomWord(c *gin.Context) {
	ctx := c.Request.Context()

	langCode, ok := c.GetQuery(QueryParamLangCode)
	if !ok {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("dictionary.delivery.Handler.getRandomWord - not found query param [lang_code]"))
		return
	}

	word, err := h.dictSvc.GetRandomWord(ctx, langCode)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("dictionary.delivery.Handler.getRandomWord: %v", err))
		return
	}

	randomWordRs := &WordRs{
		Text:          word.Text,
		LangCode:      word.LangCode,
		Pronunciation: word.Pronunciation,
	}

	c.JSON(http.StatusOK, randomWordRs)
}

package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	dictionarySvc "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	"github.com/av-ugolkov/lingua-evo/runtime"
	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

const (
	ErrMsgNotFoundWord = "Word not found"
)

const (
	QueryParamText             = "text"
	QueryParamLangCode         = "lang_code"
	QueryParamRandomWordsLimit = "limit"
	QueryParamPage             = "page"
	QueryParamPerPage          = "per_page"
	QueryParamSearch           = "search"
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

	WordDataRs struct {
		ID            *uuid.UUID `json:"id,omitempty"`
		Text          string     `json:"text,omitempty"`
		Pronunciation string     `json:"pronunciation,omitempty"`
		LangCode      string     `json:"lang_code,omitempty"`
		Creator       string     `json:"creator,omitempty"`
		CreatedAt     *time.Time `json:"created_at,omitempty"`
	}
)

type Handler struct {
	dictSvc *dictionarySvc.Service
}

func Create(r *ginext.Engine, dictSvc *dictionarySvc.Service) {
	h := newHandler(dictSvc)

	r.GET(handler.Dictionary, h.getDictionary)
	r.POST(handler.DictionaryWord, middleware.Auth(h.addWord))
	r.GET(handler.DictionaryWord, h.getWord)
	r.GET(handler.GetRandomWord, h.getRandomWords)
	r.GET(handler.WordPronunciation, middleware.Auth(h.getPronunciation))
}

func newHandler(dictSvc *dictionarySvc.Service) *Handler {
	return &Handler{
		dictSvc: dictSvc,
	}
}

func (h *Handler) getDictionary(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	langCode, err := c.GetQuery(QueryParamLangCode)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("dictionary.delivery.Handler.getDictionary: %w", err)
	}

	page, err := c.GetQueryInt(QueryParamPage)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("dictionary.delivery.Handler.getDictionary: %w", err)
	}
	itemsPerPage, err := c.GetQueryInt(QueryParamPerPage)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("dictionary.delivery.Handler.getDictionary: %w", err)
	}
	search, err := c.GetQuery(QueryParamSearch)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("dictionary.delivery.Handler.getDictionary: %w", err)
	}

	dict, countWords, err := h.dictSvc.GetDictionary(ctx, langCode, search, page, itemsPerPage)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("dictionary.delivery.Handler.getDictionary: %v", err)
	}

	wordsRs := make([]WordDataRs, 0, len(dict))
	for _, w := range dict {
		wordsRs = append(wordsRs, WordDataRs{
			ID:            &w.ID,
			Text:          w.Text,
			Pronunciation: w.Pronunciation,
			LangCode:      w.LangCode,
			Creator:       w.Creator,
			CreatedAt:     &w.CreatedAt,
		})
	}

	return http.StatusOK, gin.H{"words": wordsRs, "count_words": countWords}, nil
}

func (h *Handler) addWord(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	var data WordRq
	if err := c.Bind(&data); err != nil {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("dictionary.delivery.Handler.addWord: %v", err),
				msgerr.ErrMsgBadRequest)
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
			msgerr.New(fmt.Errorf("dictionary.delivery.Handler.addWord: %v", err),
				msgerr.ErrMsgInternal)
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

	if len(wordIDs) == 0 {
		return http.StatusNotFound, nil,
			msgerr.New(fmt.Errorf("dictionary.delivery.Handler.getWord: word not found"),
				ErrMsgNotFoundWord)
	}

	return http.StatusOK, &WordRs{
		ID: &wordIDs[0].ID,
	}, nil
}

func (h *Handler) getRandomWords(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	langCode, err := c.GetQuery(QueryParamLangCode)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("dictionary.delivery.Handler.getRandomWord: %w", err)
	}

	count, err := c.GetQueryInt(QueryParamRandomWordsLimit)
	if err != nil {
		slog.Warn(fmt.Sprintf("dictionary.delivery.Handler.getRandomWord: %v", err))
		count = 1
	}

	word, err := h.dictSvc.GetRandomWords(ctx, langCode, count)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("dictionary.delivery.Handler.getRandomWord: %v", err)
	}

	wordsRs := make([]WordRs, 0, len(word))
	for _, w := range word {
		wordsRs = append(wordsRs, WordRs{
			Text:          w.Text,
			LangCode:      w.LangCode,
			Pronunciation: w.Pronunciation,
		})
	}

	return http.StatusOK, wordsRs, nil
}

func (h *Handler) getPronunciation(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	text, err := c.GetQuery(QueryParamText)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("word.delivery.Handler.getPronunciation: %w", err)
	}

	langCode, err := c.GetQuery(QueryParamLangCode)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("word.delivery.Handler.getPronunciation: %w", err)
	}

	pronunciation, err := h.dictSvc.GetPronunciation(ctx, text, langCode)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("word.delivery.Handler.getPronunciation: %w", err)
	}

	return http.StatusOK, gin.H{"pronunciation": pronunciation}, nil
}

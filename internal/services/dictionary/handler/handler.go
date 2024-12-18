package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/fext"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	dictionarySvc "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gofiber/fiber/v2"
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

func Create(r *fiber.App, dictSvc *dictionarySvc.Service) {
	h := newHandler(dictSvc)

	r.Get(handler.Dictionary, h.getDictionary)
	r.Post(handler.DictionaryWord, middleware.Auth(h.addWord))
	r.Get(handler.DictionaryWord, h.getWord)
	r.Get(handler.GetRandomWord, h.getRandomWords)
	r.Get(handler.WordPronunciation, middleware.Auth(h.getPronunciation))
}

func newHandler(dictSvc *dictionarySvc.Service) *Handler {
	return &Handler{
		dictSvc: dictSvc,
	}
}

func (h *Handler) getDictionary(c *fiber.Ctx) error {
	ctx := c.Context()

	langCode := c.Query(QueryParamLangCode)
	if langCode == runtime.EmptyString {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("dictionary.Handler.getDictionary: not found query [%s]", QueryParamLangCode)))
	}

	page := c.QueryInt(QueryParamPage, -1)
	if page == -1 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("dictionary.Handler.getDictionary: not found query [%s]", QueryParamPage)))
	}
	itemsPerPage := c.QueryInt(QueryParamPerPage, -1)
	if itemsPerPage == -1 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("dictionary.Handler.getDictionary: not found query [%s]", QueryParamPerPage)))
	}
	search := c.Query(QueryParamSearch)

	dict, countWords, err := h.dictSvc.GetDictionary(ctx, langCode, search, page, itemsPerPage)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
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

	return c.Status(http.StatusOK).JSON(fext.D(fiber.Map{"words": wordsRs, "count_words": countWords}))
}

func (h *Handler) addWord(c *fiber.Ctx) error {
	ctx := c.Context()

	var data WordRq
	if err := c.BodyParser(&data); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
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
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgInternal))
	}

	wordRs := &WordRs{
		ID:            &words[0].ID,
		Text:          words[0].Text,
		Pronunciation: words[0].Pronunciation,
		CreatedAt:     &words[0].CreatedAt,
		UpdatedAt:     &words[0].UpdatedAt,
	}

	return c.Status(http.StatusOK).JSON(fext.D(wordRs))
}

func (h *Handler) getWord(c *fiber.Ctx) error {
	ctx := c.Context()

	text := c.Query(QueryParamText)
	if text == runtime.EmptyString {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("dictionary.Handler.getWord: not found query [%s]", QueryParamText),
			"You must specify the text of the word"))
	}

	langCode := c.Query(QueryParamLangCode)
	if langCode == runtime.EmptyString {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("dictionary.Handler.getWord: not found query [%s]", QueryParamLangCode),
			"You must specify the language code"))
	}

	wordIDs, err := h.dictSvc.GetWordsByText(ctx, []entity.DictWord{
		{
			Text:     text,
			LangCode: langCode,
		},
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgInternal))
	}

	if len(wordIDs) == 0 {
		return c.Status(http.StatusNotFound).JSON(fext.E(
			fmt.Errorf("dictionary.Handler.getWord: not found words")), ErrMsgNotFoundWord)
	}

	return c.Status(http.StatusOK).JSON(fext.D(&WordRs{
		ID: &wordIDs[0].ID,
	}))
}

func (h *Handler) getRandomWords(c *fiber.Ctx) error {
	ctx := c.Context()

	langCode := c.Query(QueryParamLangCode)
	if langCode == runtime.EmptyString {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("dictionary.Handler.getRandomWords: not found query [%s]", QueryParamLangCode),
			"You must specify the language code"))
	}

	count := c.QueryInt(QueryParamRandomWordsLimit, 1)
	word, err := h.dictSvc.GetRandomWords(ctx, langCode, count)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgInternal))
	}

	wordsRs := make([]WordRs, 0, len(word))
	for _, w := range word {
		wordsRs = append(wordsRs, WordRs{
			Text:          w.Text,
			LangCode:      w.LangCode,
			Pronunciation: w.Pronunciation,
		})
	}

	return c.Status(http.StatusOK).JSON(fext.D(wordsRs))
}

func (h *Handler) getPronunciation(c *fiber.Ctx) error {
	ctx := c.Context()

	text := c.Query(QueryParamText)
	if text == runtime.EmptyString {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("dictionary.Handler.getPronunciation: not found query [%s]", QueryParamText),
			"You must specify the text"))
	}

	langCode := c.Query(QueryParamLangCode)
	if langCode == runtime.EmptyString {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("dictionary.Handler.getPronunciation: not found query [%s]", QueryParamLangCode),
			"You must specify the language code"))
	}

	pronunciation, err := h.dictSvc.GetPronunciation(ctx, text, langCode)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgInternal))
	}

	return c.Status(http.StatusOK).JSON(fext.D(fiber.Map{"pronunciation": pronunciation}))
}

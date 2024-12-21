package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/fext"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/valid"
	dictionarySvc "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	ErrMsgNotFoundWord = "Word not found"
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

	r.Route("/dictionary", func(r fiber.Router) {
		r.Get("/", h.getDictionary)
		r.Route("/word", func(r fiber.Router) {
			r.Post("/", middleware.Auth(h.addWord))
			r.Get("/", h.getWord)
			r.Get("/random", h.getRandomWords)
			r.Get("/pronunciation", middleware.Auth(h.getPronunciation))
		})
	})
}

func newHandler(dictSvc *dictionarySvc.Service) *Handler {
	return &Handler{
		dictSvc: dictSvc,
	}
}

// getDictionary
//
// in - queries: lang_code, page, per_page, search
//
// out 200 - {words: []words: {id*, text*, pronunciation*, lang_code*, creator*, created_at*}, count_words: int}
//
// out 4..,5.. - {msg: error}
func (h *Handler) getDictionary(c *fiber.Ctx) error {
	ctx := c.Context()

	var queries struct {
		LangCode     string `query:"lang_code" validate:"required"`
		Page         int    `query:"page" validate:"required, gte=1"`
		ItemsPerPage int    `query:"per_page" validate:"required, gte=1"`
		Search       string `query:"search"`
	}
	err := c.QueryParser(&queries)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	if err := valid.Struct(queries); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	dict, countWords, err := h.dictSvc.GetDictionary(ctx,
		queries.LangCode, queries.Search, queries.Page, queries.ItemsPerPage)
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

// getWord
//
// in - queries: text, lang_code
//
// out 200 - {words: []words: {id*, text*, pronunciation*, lang_code*, creator*, created_at*}}
//
// out 4..,5.. - {msg: error}
func (h *Handler) getWord(c *fiber.Ctx) error {
	ctx := c.Context()

	var queries struct {
		Text     string `query:"text" validate:"required"`
		LangCode string `query:"lang_code" validate:"required"`
	}
	err := c.QueryParser(&queries)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	if err := valid.Struct(queries); err != nil {
		fe := err.(validator.ValidationErrors)
		for _, e := range fe {
			switch e.Namespace() {
			case "Text":
				return c.Status(http.StatusBadRequest).JSON(fext.E(e, "You must specify the text of the word"))
			case "LangCode":
				return c.Status(http.StatusBadRequest).JSON(fext.E(e, "You must specify the language code"))
			}
		}
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgValidation))
	}

	wordIDs, err := h.dictSvc.GetWordsByText(ctx, []entity.DictWord{
		{
			Text:     queries.Text,
			LangCode: queries.LangCode,
		},
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgInternal))
	}

	if len(wordIDs) == 0 {
		return c.Status(http.StatusNotFound).JSON(fext.E(
			fmt.Errorf("dictionary.Handler.getWord: not found words"), ErrMsgNotFoundWord))
	}

	return c.Status(http.StatusOK).JSON(fext.D(&WordRs{
		ID: &wordIDs[0].ID,
	}))
}

// getRandomWords
//
// in - queries: count = 1...10
//
// out 200 - {data: []words: {text*, lang_code*, pronunciation}}
//
// out 4..,5.. - {msg: error}
func (h *Handler) getRandomWords(c *fiber.Ctx) error {
	ctx := c.Context()

	var queries struct {
		Count int `query:"count" validate:"gte=1,lte=10"`
	}
	err := c.QueryParser(&queries)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}
	if queries.Count == 0 {
		queries.Count = 1
	}

	err = valid.Struct(queries)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgValidation))
	}

	word, err := h.dictSvc.GetRandomWords(ctx, queries.Count)
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

// getPronunciation (secure)
//
// in - queries: text, lang_code
//
// out 200 - {data: {pronunciation}}
//
// out 4..,5.. - {msg: error}
func (h *Handler) getPronunciation(c *fiber.Ctx) error {
	ctx := c.Context()

	var queries struct {
		Text     string `query:"text" validate:"required"`
		LangCode string `query:"lang_code" validate:"required"`
	}
	err := c.QueryParser(&queries)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	err = valid.Struct(queries)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgValidation))
	}

	pronunciation, err := h.dictSvc.GetPronunciation(ctx, queries.Text, queries.LangCode)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgInternal))
	}

	return c.Status(http.StatusOK).JSON(fext.D(fiber.Map{"pronunciation": pronunciation}))
}

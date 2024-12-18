package handler

import (
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/fext"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/router"
	"github.com/av-ugolkov/lingua-evo/internal/services/language"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gofiber/fiber/v2"
)

type (
	LanguageRs struct {
		Language string `json:"lang,omitempty"`
		Code     string `json:"code"`
	}

	Handler struct {
		langSvc *language.Service
	}
)

func Create(r *fiber.App, langSvc *language.Service) {
	h := newHandler(langSvc)

	r.Get(handler.CurrentLanguage, h.getCurrentLanguage)
	r.Get(handler.AvailableLanguages, h.getAvailableLanguages)
}

func newHandler(langSvc *language.Service) *Handler {
	return &Handler{
		langSvc: langSvc,
	}
}

func (h *Handler) getCurrentLanguage(c *fiber.Ctx) error {
	langCode := c.Cookies(router.Language, "en")
	languageRs := &LanguageRs{
		Code: langCode,
	}

	c.Cookie(&fiber.Cookie{
		Name:     router.Language,
		Value:    langCode,
		MaxAge:   0,
		Path:     "/",
		Domain:   runtime.EmptyString,
		Secure:   false,
		HTTPOnly: true,
	})

	return c.Status(http.StatusOK).JSON(fext.D(languageRs))
}

func (h *Handler) getAvailableLanguages(c *fiber.Ctx) error {
	ctx := c.Context()
	languages, err := h.langSvc.GetAvailableLanguages(ctx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	languagesRs := make([]LanguageRs, 0, len(languages))
	for _, lang := range languages {
		languagesRs = append(languagesRs, LanguageRs{
			Language: lang.Lang,
			Code:     lang.Code,
		})
	}

	return c.Status(http.StatusOK).JSON(fext.D(languagesRs))
}

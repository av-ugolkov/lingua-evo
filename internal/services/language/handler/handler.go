package handler

import (
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	"github.com/av-ugolkov/lingua-evo/internal/services/language"
	"github.com/av-ugolkov/lingua-evo/runtime"
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

func Create(r *ginext.Engine, langSvc *language.Service) {
	h := newHandler(langSvc)

	r.GET(handler.CurrentLanguage, h.getCurrentLanguage)
	r.GET(handler.AvailableLanguages, h.getAvailableLanguages)
}

func newHandler(langSvc *language.Service) *Handler {
	return &Handler{
		langSvc: langSvc,
	}
}

func (h *Handler) getCurrentLanguage(c *ginext.Context) (int, any, error) {
	langCode, err := c.Cookie(ginext.Language)
	if err != nil {
		langCode = runtime.GetLanguage("en")
	}
	languageRs := &LanguageRs{
		Code: langCode,
	}

	c.SetCookie(ginext.Language, languageRs.Code, 0, "/", runtime.EmptyString, false, true)
	return http.StatusOK, languageRs, nil
}

func (h *Handler) getAvailableLanguages(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()
	languages, err := h.langSvc.GetAvailableLanguages(ctx)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	languagesRs := make([]LanguageRs, 0, len(languages))
	for _, lang := range languages {
		languagesRs = append(languagesRs, LanguageRs{
			Language: lang.Lang,
			Code:     lang.Code,
		})
	}

	return http.StatusOK, languagesRs, nil
}

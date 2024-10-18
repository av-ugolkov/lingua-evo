package handler

import (
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/delivery/handler/gin"
	"github.com/av-ugolkov/lingua-evo/internal/services/language"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gin-gonic/gin"
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

func Create(r *gin.Engine, langSvc *language.Service) {
	h := newHandler(langSvc)
	h.register(r)
}

func newHandler(langSvc *language.Service) *Handler {
	return &Handler{
		langSvc: langSvc,
	}
}

func (h *Handler) register(r *gin.Engine) {
	r.GET(handler.CurrentLanguage, h.getCurrentLanguage)
	r.GET(handler.AvailableLanguages, h.getAvailableLanguages)
}

func (h *Handler) getCurrentLanguage(c *gin.Context) {
	langCode, err := c.Cookie(ginExt.Language)
	if err != nil {
		langCode = runtime.GetLanguage("en")
	}
	languageRs := &LanguageRs{
		Code: langCode,
	}

	c.SetCookie(ginExt.Language, languageRs.Code, 0, "/", "", false, true)
	c.JSON(http.StatusOK, languageRs)
}

func (h *Handler) getAvailableLanguages(c *gin.Context) {
	ctx := c.Request.Context()
	languages, err := h.langSvc.GetAvailableLanguages(ctx)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError, err)
		return
	}

	languagesRs := make([]LanguageRs, 0, len(languages))
	for _, lang := range languages {
		languagesRs = append(languagesRs, LanguageRs{
			Language: lang.Lang,
			Code:     lang.Code,
		})
	}

	c.JSON(http.StatusOK, languagesRs)
}

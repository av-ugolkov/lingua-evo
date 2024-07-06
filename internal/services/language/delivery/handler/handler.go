package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/pkg/http/gin_extension"
	"github.com/av-ugolkov/lingua-evo/internal/services/language"
	"github.com/av-ugolkov/lingua-evo/internal/services/language/dto"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	langSvc *language.Service
}

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
	r.GET(delivery.CurrentLanguage, h.getCurrentLanguage)
	r.GET(delivery.AvailableLanguages, h.getAvailableLanguages)
}

func (h *Handler) getCurrentLanguage(c *gin.Context) {
	langCode, err := c.Cookie(ginExt.Language)
	if err != nil {
		langCode = runtime.GetLanguage("en")
	}
	languageRs := &dto.LanguageRs{
		Code: langCode,
	}

	c.SetCookie(ginExt.Language, languageRs.Code, 0, "/", "", false, true)
	c.JSON(http.StatusOK, languageRs)
}

func (h *Handler) getAvailableLanguages(c *gin.Context) {
	ctx := c.Request.Context()
	languages, err := h.langSvc.GetAvailableLanguages(ctx)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("lingua.language.delivery.Handler.getAvailableLanguages: %v", err))
		return
	}

	languagesRs := make([]dto.LanguageRs, 0, len(languages))
	for _, lang := range languages {
		languagesRs = append(languagesRs, dto.LanguageRs{
			Language: lang.Lang,
			Code:     lang.Code,
		})
	}

	c.JSON(http.StatusOK, languagesRs)
}

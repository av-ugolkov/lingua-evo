package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"lingua-evo/internal/services/language/service"
	"lingua-evo/pkg/http/exchange"
)

const (
	getCurrentLanguage    = "/get_current_languages"
	getAvailableLanguages = "/get_languages"
)

type (
	LanguageRs struct {
		Language string `json:"language"`
		Code     string `json:"code"`
	}

	Handler struct {
		langSvc *service.LanguageSvc
	}
)

func Create(r *mux.Router, langSvc *service.LanguageSvc) {
	h := newHandler(langSvc)
	h.register(r)
}

func newHandler(langSvc *service.LanguageSvc) *Handler {
	return &Handler{
		langSvc: langSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(getCurrentLanguage, h.getCurrentLanguage).Methods(http.MethodGet)
	r.HandleFunc(getAvailableLanguages, h.getAvailableLanguages).Methods(http.MethodGet)
}

func (h *Handler) getCurrentLanguage(w http.ResponseWriter, r *http.Request) {
	ex := exchange.NewExchanger(w, r)

	languageRs := &LanguageRs{
		Code: ex.GetCookieLanguageOrDefault(),
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SetCookieLanguage(languageRs.Code)
	ex.SendData(http.StatusOK, languageRs)
}

func (h *Handler) getAvailableLanguages(w http.ResponseWriter, r *http.Request) {
	ex := exchange.NewExchanger(w, r)

	languages, err := h.langSvc.GetAvailableLanguages(r.Context())
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("lingua.language.delivery.Handler.getAvailableLanguages: %v", err))
		return
	}

	languagesRs := make([]LanguageRs, 0, len(languages))
	for _, lang := range languages {
		languagesRs = append(languagesRs, LanguageRs{
			Language: lang.Lang,
			Code:     lang.Code,
		})
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, languagesRs)
}
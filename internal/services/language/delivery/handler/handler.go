package handler

import (
	"encoding/json"
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

	type Language struct {
		Language string `json:"language"`
		Code     string `json:"code"`
	}
	lang := Language{}
	lang.Code = ex.GetCookieLanguageOrDefault()

	b, err := json.Marshal(lang)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("lingua.language.delivery.Handler.getCurrentLanguage - marshal: %v", err))
		return
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SetCookieLanguage(lang.Code)
	ex.SendData(http.StatusOK, b)
}

func (h *Handler) getAvailableLanguages(w http.ResponseWriter, r *http.Request) {
	ex := exchange.NewExchanger(w, r)

	languages, err := h.langSvc.GetAvailableLanguages(r.Context())
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("lingua.language.delivery.Handler.getAvailableLanguages: %v", err))
		return
	}

	jsonLanguages, err := json.Marshal(languages)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("lingua.language.delivery.Handler.getAvailableLanguages - marshal: %v", err))
		return
	}
	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, jsonLanguages)
}

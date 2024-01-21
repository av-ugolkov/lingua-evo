package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"lingua-evo/internal/services/lingua/language/service"
	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/http/handler/common"
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
	handler := newHandler(langSvc)
	handler.register(r)
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
	handler := handler.NewHandler(w, r)

	type Language struct {
		Language string `json:"language"`
		Code     string `json:"code"`
	}
	lang := Language{}
	lang.Code = handler.GetCookieLanguageOrDefault()

	b, err := json.Marshal(lang)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("lingua.language.delivery.Handler.getCurrentLanguage - marshal: %v", err))
		return
	}

	handler.SetContentType(common.ContentTypeJSON)
	handler.SetCookieLanguage(lang.Code)
	handler.SendData(http.StatusOK, b)
}

func (h *Handler) getAvailableLanguages(w http.ResponseWriter, r *http.Request) {
	handler := handler.NewHandler(w, r)

	languages, err := h.langSvc.GetAvailableLanguages(r.Context())
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("lingua.language.delivery.Handler.getAvailableLanguages: %v", err))
		return
	}

	jsonLanguages, err := json.Marshal(languages)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("lingua.language.delivery.Handler.getAvailableLanguages - marshal: %v", err))
		return
	}
	handler.SetContentType(common.ContentTypeJSON)
	handler.SendData(http.StatusOK, jsonLanguages)
}

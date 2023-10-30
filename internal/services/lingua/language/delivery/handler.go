package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"lingua-evo/internal/services/lingua/language/service"
)

const (
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
	r.HandleFunc(getAvailableLanguages, h.getAvailableLanguages).Methods(http.MethodGet)
}

func (h *Handler) getAvailableLanguages(w http.ResponseWriter, r *http.Request) {
	languages, err := h.langSvc.GetAvailableLanguages(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	jsonLanguages, err := json.Marshal(languages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	_, _ = w.Write(jsonLanguages)
}

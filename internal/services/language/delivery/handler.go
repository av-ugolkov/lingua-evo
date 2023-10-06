package delivery

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"lingua-evo/internal/services/language/entity"
)

const (
	getLanguages = "/get_languages"
)

type (
	langSvc interface {
		GetLanguages(context.Context) ([]*entity.Language, error)
	}

	Handler struct {
		langSvc langSvc
	}
)

func Create(r *mux.Router, langSvc langSvc) {
	handler := newHandler(langSvc)
	handler.register(r)
}

func newHandler(langSvc langSvc) *Handler {
	return &Handler{
		langSvc: langSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(getLanguages, h.getLanguages).Methods(http.MethodGet)
}

func (h *Handler) getLanguages(w http.ResponseWriter, r *http.Request) {
	languages, err := h.langSvc.GetLanguages(r.Context())
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

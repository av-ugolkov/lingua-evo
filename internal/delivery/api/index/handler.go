package index

import (
	"html/template"
	"lingua-evo/internal/config"
	"lingua-evo/internal/service"
	"lingua-evo/pkg/logging"
	"lingua-evo/pkg/tools/view"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	url      = "/"
	indexURL = "/index"

	indexPagePath = "view/index.html"
)

type Handler struct {
	logger *logging.Logger
	lingua *service.Lingua
}

func Create(log *logging.Logger, ling *service.Lingua, r *httprouter.Router) {
	handler := newHandler(log, ling)
	handler.register(r)
}

func newHandler(logger *logging.Logger, lingua *service.Lingua) *Handler {
	return &Handler{
		logger: logger,
		lingua: lingua,
	}
}

func (h *Handler) register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, url, h.getIndex)
	router.HandlerFunc(http.MethodGet, indexURL, h.getIndex)
}

func (h *Handler) getIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(view.GetPathFile(indexPagePath))
	if err != nil {
		h.logger.Errorf("index.get.ParseFiles: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data := struct {
		RootPath string
	}{
		RootPath: config.GetConfig().Front.Root,
	}

	err = t.Execute(w, data)
	if err != nil {
		h.logger.Errorf("index.get.Execute: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

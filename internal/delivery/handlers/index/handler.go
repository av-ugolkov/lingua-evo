package index

import (
	"net/http"

	"lingua-evo/internal/services"
	"lingua-evo/pkg/logging"
	staticFiles "lingua-evo/static"

	"github.com/julienschmidt/httprouter"
)

const (
	url      = "/"
	indexURL = "/index"

	indexPagePath = "web/index.html"
)

type Handler struct {
	logger *logging.Logger
	lingua *services.Lingua
}

func Create(log *logging.Logger, ling *services.Lingua, r *httprouter.Router) {
	handler := newHandler(log, ling)
	handler.register(r)
}

func newHandler(logger *logging.Logger, lingua *services.Lingua) *Handler {
	return &Handler{
		logger: logger,
		lingua: lingua,
	}
}

func (h *Handler) register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, url, h.get)
	router.HandlerFunc(http.MethodGet, indexURL, h.get)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	file, err := staticFiles.OpenFile(indexPagePath)
	if err != nil {
		h.logger.Errorf("index.get.OpenFile: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(file))
}

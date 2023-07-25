package index

import (
	"lingua-evo/internal/service"
	"lingua-evo/pkg/logging"
	templates "lingua-evo/web/static"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	url      = "/"
	indexURL = "/index"

	indexPagePath = "index.html"
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
	router.HandlerFunc(http.MethodGet, url, h.get)
	router.HandlerFunc(http.MethodGet, indexURL, h.get)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	file, err := templates.OpenFile(indexPagePath)
	if err != nil {
		h.logger.Errorf("index.get.OpenFile: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write([]byte(file))
	w.WriteHeader(http.StatusOK)
}

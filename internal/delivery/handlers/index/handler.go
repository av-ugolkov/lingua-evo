package index

import (
	"lingua-evo/internal/service"
	"lingua-evo/pkg/logging"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

const (
	url      = "/"
	indexURL = "/index"

	indexPagePath = "./web/static/index.html"
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
	file, err := os.ReadFile(indexPagePath)
	if err != nil {
		h.logger.Errorf("index.get.ParseFiles: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write([]byte(file))
	w.WriteHeader(http.StatusOK)
}

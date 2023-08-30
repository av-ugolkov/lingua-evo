package index

import (
	"log/slog"
	"net/http"

	staticFiles "lingua-evo/static"

	"github.com/julienschmidt/httprouter"
)

const (
	url      = "/"
	indexURL = "/index"

	indexPagePath = "web/index.html"
)

type Handler struct{}

func Create(r *httprouter.Router) {
	handler := newHandler()
	handler.register(r)
}

func newHandler() *Handler {
	return &Handler{}
}

func (h *Handler) register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, url, h.get)
	router.HandlerFunc(http.MethodGet, indexURL, h.get)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	file, err := staticFiles.OpenFile(indexPagePath)
	if err != nil {
		slog.Error("index.get.OpenFile: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(file))
}

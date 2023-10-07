package index

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"

	staticFiles "lingua-evo/static"
)

const (
	mainURL  = "/"
	indexURL = "/index"

	indexPagePath = "web/index.html"
)

type Handler struct{}

func Create(r *mux.Router) {
	handler := newHandler()
	handler.register(r)
}

func newHandler() *Handler {
	return &Handler{}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(mainURL, h.get).Methods(http.MethodGet)
	r.HandleFunc(indexURL, h.get).Methods(http.MethodGet)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	file, err := staticFiles.OpenFile(indexPagePath)
	if err != nil {
		slog.Error(fmt.Errorf("index.get.OpenFile: %v", err).Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(file))
}

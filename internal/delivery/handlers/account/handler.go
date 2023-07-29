package account

import (
	"lingua-evo/internal/services"
	"lingua-evo/pkg/logging"
	staticFiles "lingua-evo/static"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	account = "/account"

	accountPagePath = "web/account/account.html"
)

type Handler struct {
	logger *logging.Logger
}

func Create(log *logging.Logger, _ *services.Lingua, r *httprouter.Router) {
	handler := newHandler(log)
	handler.register(r)
}

func newHandler(logger *logging.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, account, h.account)
}

func (h *Handler) account(w http.ResponseWriter, _ *http.Request) {
	file, err := staticFiles.OpenFile(accountPagePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = w.Write(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

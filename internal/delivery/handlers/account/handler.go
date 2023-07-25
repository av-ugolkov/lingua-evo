package account

import (
	"lingua-evo/internal/service"
	"lingua-evo/pkg/logging"
	templates "lingua-evo/web/static"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	account = "/account"

	accountPagePath = "account/account.html"
)

type Handler struct {
	logger *logging.Logger
}

func Create(log *logging.Logger, _ *service.Lingua, r *httprouter.Router) {
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
	file, err := templates.OpenFile(accountPagePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = w.Write(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

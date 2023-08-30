package delivery

import (
	"net/http"

	staticFiles "lingua-evo/static"

	"github.com/julienschmidt/httprouter"
)

const (
	account = "/account"

	accountPagePath = "web/account/account.html"
)

type Handler struct {
}

func Create(r *httprouter.Router) {
	handler := newHandler()
	handler.register(r)
}

func newHandler() *Handler {
	return &Handler{}
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

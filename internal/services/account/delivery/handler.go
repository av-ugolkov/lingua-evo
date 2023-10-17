package delivery

import (
	"net/http"

	staticFiles "lingua-evo"

	"github.com/gorilla/mux"
)

const (
	account = "/account"

	accountPagePath = "account/account.html"
)

type Handler struct {
}

func Create(r *mux.Router) {
	handler := newHandler()
	handler.register(r)
}

func newHandler() *Handler {
	return &Handler{}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(account, h.account).Methods(http.MethodGet)
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

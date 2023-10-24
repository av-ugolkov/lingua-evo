package sign_up

import (
	"fmt"
	"log/slog"
	"net/http"

	staticFiles "lingua-evo"

	"github.com/gorilla/mux"
)

const (
	signUpURL = "/signup"

	signupPage = "website/sign_up/signup.html"
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
	r.HandleFunc(signUpURL, h.get).Methods(http.MethodGet)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	file, err := staticFiles.OpenFile(signupPage)
	if err != nil {
		slog.Error(fmt.Errorf("sign_up.get.OpenFile: %v", err).Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(file))
}

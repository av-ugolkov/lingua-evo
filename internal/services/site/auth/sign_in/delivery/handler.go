package delivery

import (
	"fmt"
	"log/slog"
	"net/http"

	"lingua-evo/internal/services/user/service"
	"lingua-evo/pkg/files"

	"github.com/gorilla/mux"
)

const (
	signInURL = "/signin"

	signInPage = "website/sign_in/signin.html"
)

type (
	Handler struct {
		userSvc *service.UserSvc
	}
)

func Create(r *mux.Router, userSvc *service.UserSvc) {
	handler := newHandler(userSvc)
	handler.register(r)
}

func newHandler(userSvc *service.UserSvc) *Handler {
	return &Handler{
		userSvc: userSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(signInURL, h.get).Methods(http.MethodGet)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	file, err := files.OpenFile(signInPage)
	if err != nil {
		slog.Error(fmt.Errorf("sign_in.get.OpenFile: %v", err).Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(file))
}

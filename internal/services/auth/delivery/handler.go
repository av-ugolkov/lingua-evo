package delivery

import (
	"net/http"

	"lingua-evo/internal/services/user/service"

	"github.com/gorilla/mux"
)

const (
	createSession = "/signin"
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
	r.HandleFunc(createSession, h.createSession).Methods(http.MethodPost)
}

func (h *Handler) createSession(w http.ResponseWriter, r *http.Request) {

}

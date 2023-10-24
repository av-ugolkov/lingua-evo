package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"lingua-evo/internal/services/auth/dto"
	"lingua-evo/internal/services/auth/service"
	"lingua-evo/pkg/tools"

	"github.com/gorilla/mux"
)

const (
	createSession = "/signin"
)

type (
	Handler struct {
		authSvc *service.AuthSvc
	}
)

func Create(r *mux.Router, authSvc *service.AuthSvc) {
	handler := newHandler(authSvc)
	handler.register(r)
}

func newHandler(authSvc *service.AuthSvc) *Handler {
	return &Handler{
		authSvc: authSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(createSession, h.createSession).Methods(http.MethodPost)
}

func (h *Handler) createSession(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	var data dto.CreateSessionRq

	err := tools.CheckBody(w, r, &data)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.createSession - check body: %v", err))
		return
	}

	ctx := r.Context()
	tokens, err := h.authSvc.CreateSession(ctx, &data)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.createSession - create session: %v", err))
		return
	}

	b, err := json.Marshal(&dto.CreateSessionRs{
		JWT:          tokens.JWT,
		RefreshToken: tokens.RefreshToken,
	})
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.createSession - marshal: %v", err))
		return
	}

	_, _ = w.Write(b)
}

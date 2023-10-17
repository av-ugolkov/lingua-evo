package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"lingua-evo/internal/services/user/dto"
	"lingua-evo/internal/services/user/service"
	"lingua-evo/internal/tools"

	"github.com/gorilla/mux"
)

const (
	userGetIDByName  = "/user/get_id_by_name"
	userGetIDByEmail = "/user/get_id_by_email"
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
	r.HandleFunc(userGetIDByName, h.getIDByName).Methods(http.MethodPost)
	r.HandleFunc(userGetIDByEmail, h.getIDByEmail).Methods(http.MethodPost)
}

func (h *Handler) getIDByName(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	var data dto.GetIDRq

	err := tools.CheckBody(w, r, &data)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByName - check body: %v", err))
		return
	}

	ctx := r.Context()
	id, err := h.userSvc.GetIDByName(ctx, data.Value)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByName: %v", err))
		return
	}

	userID := dto.UserIDRs{
		ID: id,
	}
	b, err := json.Marshal(userID)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByName - marshal: %v", err))
		return
	}
	_, _ = w.Write(b)
}

func (h *Handler) getIDByEmail(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	var data dto.GetIDRq

	err := tools.CheckBody(w, r, &data)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByEmail - check body: %v", err))
		return
	}

	ctx := r.Context()
	id, err := h.userSvc.GetIDByEmail(ctx, data.Value)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByEmail: %v", err))
		return
	}

	userID := dto.UserIDRs{
		ID: id,
	}
	b, err := json.Marshal(userID)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByEmail - marshal: %v", err))
		return
	}
	_, _ = w.Write(b)
}

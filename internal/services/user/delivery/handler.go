package delivery

import (
	"context"
	"fmt"
	"lingua-evo/internal/services/user/dto"
	"lingua-evo/internal/tools"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	userGetIDByName  = "/user/get_id_by_name"
	userGetIDByEmail = "/user/get_id_by_email"
)

type (
	userSvc interface {
		GetIDByName(ctx context.Context, name string) (uuid.UUID, error)
		GetIDByEmail(ctx context.Context, email string) (uuid.UUID, error)
	}

	Handler struct {
		userSvc userSvc
	}
)

func Create(r *mux.Router, userSvc userSvc) {
	handler := newHandler(userSvc)
	handler.register(r)
}

func newHandler(userSvc userSvc) *Handler {
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
	userID, err := h.userSvc.GetIDByName(ctx, data.Value)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByName: %v", err))
		return
	}

	_, _ = w.Write([]byte(userID.String()))
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
	userID, err := h.userSvc.GetIDByEmail(ctx, data.Value)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByEmail: %v", err))
		return
	}

	_, _ = w.Write([]byte(userID.String()))
}

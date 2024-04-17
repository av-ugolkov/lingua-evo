package handler

import (
	"context"
	"fmt"
	"github.com/av-ugolkov/lingua-evo/internal/services/user/model"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/http/exchange"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/services/user"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gorilla/mux"
)

type Handler struct {
	userSvc *user.Service
}

func Create(r *mux.Router, userSvc *user.Service) {
	h := newHandler(userSvc)
	h.register(r)
}

func newHandler(userSvc *user.Service) *Handler {
	return &Handler{
		userSvc: userSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(delivery.UserByID, middleware.Auth(h.getUserByID)).Methods(http.MethodGet)
}

func (h *Handler) getUserByID(ctx context.Context, ex *exchange.Exchanger) {
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ex.SendError(http.StatusUnauthorized, fmt.Errorf("user.delivery.Handler.getUserByID - unauthorized: %v", err))
		return
	}
	userData, err := h.userSvc.GetUserByID(ctx, userID)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getUserByID: %v", err))
		return
	}

	userRs := &model.UserRs{
		ID:    userData.ID,
		Name:  userData.Name,
		Email: userData.Email,
		Role:  userData.Role,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, userRs)
}

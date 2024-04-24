package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/http/exchange"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/services/user"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/user"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type (
	CreateUserRq struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Code     int    `json:"code"`
	}

	GetValueRq struct {
		Value string `json:"value"`
	}

	CreateUserRs struct {
		UserID uuid.UUID `json:"user_id"`
	}

	UserRs struct {
		ID    uuid.UUID    `json:"id"`
		Name  string       `json:"name"`
		Email string       `json:"email"`
		Role  runtime.Role `json:"role"`
	}
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
	r.HandleFunc(delivery.SignUp, h.signUp).Methods(http.MethodPost)
	r.HandleFunc(delivery.UserByID, middleware.Auth(h.getUserByID)).Methods(http.MethodGet)
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	ex := exchange.NewExchanger(w, r)
	var data CreateUserRq
	err := ex.CheckBody(&data)
	if err != nil {
		ex.SendError(http.StatusBadRequest, fmt.Errorf("user.delivery.Handler.createAccount - check body: %v", err))
		return
	}

	uid, err := h.userSvc.SignUp(ex.Context(), entity.UserData{
		ID:       uuid.New(),
		Name:     data.Username,
		Password: data.Password,
		Email:    data.Email,
		Role:     runtime.User,
		Code:     data.Code,
	})
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.createAccount - create user: %v", err))
		return
	}

	createUserRs := &CreateUserRs{
		UserID: uid,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusCreated, createUserRs)
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

	userRs := &UserRs{
		ID:    userData.ID,
		Name:  userData.Name,
		Email: userData.Email,
		Role:  userData.Role,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusOK, userRs)
}

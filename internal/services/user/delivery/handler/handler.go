package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	entity "lingua-evo/internal/services/user"
	"lingua-evo/internal/services/user/service"
	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/http/handler/common"
	"lingua-evo/pkg/middleware"
	"lingua-evo/pkg/utils"
	"lingua-evo/runtime"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	createAccount  = "/signup"
	getUserByID    = "/user/get_by_id"
	getUserByName  = "/user/get_by_name"
	getUserByEmail = "/user/get_by_email"
)

type (
	CreateUserRq struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	CreateUserRs struct {
		UserID uuid.UUID `json:"user_id"`
	}

	GetValueRq struct {
		Value string `json:"value"`
	}

	UserRs struct {
		ID    uuid.UUID    `json:"id"`
		Name  string       `json:"name"`
		Email string       `json:"email"`
		Role  runtime.Role `json:"role"`
	}

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
	r.HandleFunc(createAccount, h.createAccount).Methods(http.MethodPost)
	r.HandleFunc(getUserByID, middleware.Auth(h.getUserByID)).Methods(http.MethodGet)
	r.HandleFunc(getUserByName, middleware.Auth(h.getUserByName)).Methods(http.MethodPost)
	r.HandleFunc(getUserByEmail, middleware.Auth(h.getUserByEmail)).Methods(http.MethodPost)
}

func (h *Handler) createAccount(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	handler := handler.NewHandler(w, r)
	var data CreateUserRq
	err := handler.CheckBody(&data)
	if err != nil {
		handler.SendError(http.StatusBadRequest, fmt.Errorf("user.delivery.Handler.createAccount - check body: %v", err))
		return
	}

	if err := h.validateEmail(r.Context(), data.Email); err != nil {
		handler.SendError(http.StatusConflict, fmt.Errorf("user.delivery.Handler.createAccount - validateEmail: %v", err))
		return
	}

	if err := h.validateUsername(r.Context(), data.Username); err != nil {
		handler.SendError(http.StatusConflict, fmt.Errorf("user.delivery.Handler.createAccount - validateUsername: %v", err))
		return
	}

	if err := validatePassword(data.Password); err != nil {
		handler.SendError(http.StatusConflict, fmt.Errorf("user.delivery.Handler.createAccount - validatePassword: %v", err))
		return
	}

	hashPassword, err := utils.HashPassword(data.Password)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.createAccount - hashPassword: %v", err))
		return
	}

	data.Password = hashPassword

	uid, err := h.userSvc.CreateUser(r.Context(), data.Username, data.Password, data.Email, runtime.User)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.createAccount - create user: %v", err))
		return
	}

	b, err := json.Marshal(&CreateUserRs{
		UserID: uid,
	})
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.createAccount - marshal: %v", err))
		return
	}

	handler.SetContentType(common.ContentTypeJSON)
	handler.SendData(http.StatusCreated, b)
}

func (h *Handler) getUserByID(w http.ResponseWriter, r *http.Request) {
	handler := handler.NewHandler(w, r)

	ctx := r.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		handler.SendError(http.StatusUnauthorized, fmt.Errorf("user.delivery.Handler.getUserByID - unauthorized: %v", err))
		return
	}
	user, err := h.userSvc.GetUserByID(ctx, userID)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getUserByID: %v", err))
		return
	}

	userRs := UserRs{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
	b, err := json.Marshal(userRs)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getUserByID - marshal: %v", err))
		return
	}
	handler.SetContentType(common.ContentTypeJSON)
	handler.SendData(http.StatusOK, b)
}

func (h *Handler) getUserByName(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	handler := handler.NewHandler(w, r)
	var data GetValueRq

	err := handler.CheckBody(&data)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByName - check body: %v", err))
		return
	}

	ctx := r.Context()
	user, err := h.userSvc.GetUserByName(ctx, data.Value)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByName: %v", err))
		return
	}

	userID := UserRs{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
	b, err := json.Marshal(userID)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByName - marshal: %v", err))
		return
	}
	handler.SetContentType(common.ContentTypeJSON)
	handler.SendData(http.StatusOK, b)
}

func (h *Handler) getUserByEmail(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()
	handler := handler.NewHandler(w, r)

	var data GetValueRq

	if err := handler.CheckBody(&data); err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByEmail - check body: %v", err))
		return
	}

	ctx := r.Context()
	user, err := h.userSvc.GetUserByEmail(ctx, data.Value)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByEmail: %v", err))
		return
	}

	userID := UserRs{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
	b, err := json.Marshal(userID)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByEmail - marshal: %v", err))
		return
	}
	handler.SetContentType(common.ContentTypeJSON)
	handler.SendData(http.StatusOK, b)
}

func (h *Handler) validateEmail(ctx context.Context, email string) error {
	if !utils.IsEmailValid(email) {
		return entity.ErrEmailNotCorrect
	}

	user, err := h.userSvc.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	} else if user.ID == uuid.Nil && err == nil {
		return entity.ErrItIsAdmin
	} else if user.ID != uuid.Nil {
		return entity.ErrEmailBusy
	}

	return nil
}

func (h *Handler) validateUsername(ctx context.Context, username string) error {
	if len(username) <= entity.UsernameLen {
		return entity.ErrUsernameLen
	}

	user, err := h.userSvc.GetUserByName(ctx, username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	} else if user.ID == uuid.Nil && err == nil {
		return entity.ErrItIsAdmin
	} else if user.ID != uuid.Nil {
		return entity.ErrUsernameBusy
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < entity.MinPasswordLen {
		return entity.ErrPasswordLen
	}

	if !utils.IsPasswordValid(password) {
		return entity.ErrPasswordDifficult
	}

	return nil
}

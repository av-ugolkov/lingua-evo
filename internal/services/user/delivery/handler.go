package delivery

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"lingua-evo/internal/services/user/dto"
	"lingua-evo/internal/services/user/entity"
	"lingua-evo/internal/services/user/service"
	"lingua-evo/pkg/tools"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	createAccount  = "/signup"
	getUserByName  = "/user/get_by_name"
	getUserByEmail = "/user/get_by_email"
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
	r.HandleFunc(createAccount, h.createAccount).Methods(http.MethodPost)
	r.HandleFunc(getUserByName, h.getUserByName).Methods(http.MethodGet)
	r.HandleFunc(getUserByEmail, h.getUserByEmail).Methods(http.MethodGet)
}

func (h *Handler) createAccount(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	var data dto.CreateUserRq
	err := tools.CheckBody(w, r, &data)
	if err != nil {
		tools.SendError(w, http.StatusBadRequest, fmt.Errorf("user.delivery.Handler.createAccount - check body: %v", err))
		return
	}

	if err := h.validateEmail(r.Context(), data.Email); err != nil {
		tools.SendError(w, http.StatusConflict, fmt.Errorf("user.delivery.Handler.createAccount - validateEmail: %v", err))
		return
	}

	if err := h.validateUsername(r.Context(), data.Username); err != nil {
		tools.SendError(w, http.StatusConflict, fmt.Errorf("user.delivery.Handler.createAccount - validateUsername: %v", err))
		return
	}

	if err := validatePassword(data.Password); err != nil {
		tools.SendError(w, http.StatusConflict, fmt.Errorf("user.delivery.Handler.createAccount - validatePassword: %v", err))
		return
	}

	hashPassword, err := tools.HashPassword(data.Password)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.createAccount - hashPassword: %v", err))
		return
	}

	data.Password = hashPassword

	uid, err := h.userSvc.CreateUser(r.Context(), &data)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.createAccount - create user: %v", err))
		return
	}

	b, err := json.Marshal(&dto.CreateUserRs{
		UserID: uid,
	})
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.createAccount - marshal: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(b)
}

func (h *Handler) getUserByName(w http.ResponseWriter, r *http.Request) {
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
	user, err := h.userSvc.GetUserByName(ctx, data.Value)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByName: %v", err))
		return
	}

	userID := dto.UserIDRs{
		ID: user.ID,
	}
	b, err := json.Marshal(userID)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByName - marshal: %v", err))
		return
	}
	_, _ = w.Write(b)
}

func (h *Handler) getUserByEmail(w http.ResponseWriter, r *http.Request) {
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
	user, err := h.userSvc.GetUserByEmail(ctx, data.Value)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByEmail: %v", err))
		return
	}

	userID := dto.UserIDRs{
		ID: user.ID,
	}
	b, err := json.Marshal(userID)
	if err != nil {
		tools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByEmail - marshal: %v", err))
		return
	}
	_, _ = w.Write(b)
}

func (h *Handler) validateEmail(ctx context.Context, email string) error {
	if !tools.IsEmailValid(email) {
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

	if !tools.IsPasswordValid(password) {
		return entity.ErrPasswordDifficult
	}

	return nil
}

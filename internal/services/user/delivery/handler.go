package delivery

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"lingua-evo/internal/services/auth/entity"
	"lingua-evo/internal/services/user/dto"
	"lingua-evo/internal/services/user/service"
	"lingua-evo/pkg/tools"

	httpTools "lingua-evo/internal/tools"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

const (
	createAccount    = "/signup"
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
	r.HandleFunc(createAccount, h.createAccount).Methods(http.MethodPost)
	r.HandleFunc(userGetIDByName, h.getIDByName).Methods(http.MethodPost)
	r.HandleFunc(userGetIDByEmail, h.getIDByEmail).Methods(http.MethodPost)
}

func (h *Handler) createAccount(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	var data dto.CreateUserRq
	err := httpTools.CheckBody(w, r, &data)
	if err != nil {
		httpTools.SendError(w, http.StatusBadRequest, fmt.Errorf("user.delivery.Handler.createAccount - check body: %v", err))
		return
	}

	if err := h.validateEmail(r.Context(), data.Email); err != nil {
		httpTools.SendError(w, http.StatusConflict, fmt.Errorf("user.delivery.Handler.createAccount - validateEmail: %v", err))
		return
	}

	if err := h.validateUsername(r.Context(), data.Username); err != nil {
		httpTools.SendError(w, http.StatusConflict, fmt.Errorf("user.delivery.Handler.createAccount - validateUsername: %v", err))
		return
	}

	if err := validatePassword(data.Password); err != nil {
		httpTools.SendError(w, http.StatusConflict, fmt.Errorf("user.delivery.Handler.createAccount - validatePassword: %v", err))
		return
	}

	hashPassword, err := hashPassword(data.Password)
	if err != nil {
		httpTools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.createAccount - hashPassword: %v", err))
		return
	}

	data.Password = hashPassword

	uid, err := h.userSvc.CreateUser(r.Context(), &data)
	if err != nil {
		httpTools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.createAccount - create user: %v", err))
		return
	}

	b, err := json.Marshal(&dto.CreateUserRs{
		UserID: uid,
	})
	if err != nil {
		httpTools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.createAccount - marshal: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(b)
}

func (h *Handler) getIDByName(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	var data dto.GetIDRq

	err := httpTools.CheckBody(w, r, &data)
	if err != nil {
		httpTools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByName - check body: %v", err))
		return
	}

	ctx := r.Context()
	id, err := h.userSvc.GetIDByName(ctx, data.Value)
	if err != nil {
		httpTools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByName: %v", err))
		return
	}

	userID := dto.UserIDRs{
		ID: id,
	}
	b, err := json.Marshal(userID)
	if err != nil {
		httpTools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByName - marshal: %v", err))
		return
	}
	_, _ = w.Write(b)
}

func (h *Handler) getIDByEmail(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	var data dto.GetIDRq

	err := httpTools.CheckBody(w, r, &data)
	if err != nil {
		httpTools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByEmail - check body: %v", err))
		return
	}

	ctx := r.Context()
	id, err := h.userSvc.GetIDByEmail(ctx, data.Value)
	if err != nil {
		httpTools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByEmail: %v", err))
		return
	}

	userID := dto.UserIDRs{
		ID: id,
	}
	b, err := json.Marshal(userID)
	if err != nil {
		httpTools.SendError(w, http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.getIDByEmail - marshal: %v", err))
		return
	}
	_, _ = w.Write(b)
}

func (h *Handler) validateEmail(ctx context.Context, email string) error {
	if !tools.IsEmailValid(email) {
		return entity.ErrEmailNotCorrect
	}

	uid, err := h.userSvc.GetIDByEmail(ctx, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	} else if uid == uuid.Nil && err == nil {
		return entity.ErrItIsAdmin
	} else if uid != uuid.Nil {
		return entity.ErrEmailBusy
	}

	return nil
}

func (h *Handler) validateUsername(ctx context.Context, username string) error {
	if len(username) <= entity.UsernameLen {
		return entity.ErrUsernameLen
	}

	uid, err := h.userSvc.GetIDByName(ctx, username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	} else if uid == uuid.Nil && err == nil {
		return entity.ErrItIsAdmin
	} else if uid != uuid.Nil {
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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), entity.PasswordSolt)
	return string(bytes), err
}

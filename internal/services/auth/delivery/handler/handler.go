package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"lingua-evo/internal/config"
	"lingua-evo/internal/services/auth/service"
	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/http/handler/common"
	"lingua-evo/pkg/middleware"
	"lingua-evo/runtime"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	login   = "/auth/login"
	refresh = "/auth/refresh"
	logout  = "/auth/logout"
)

type (
	CreateSessionRq struct {
		User        string `json:"user"`
		Password    string `json:"password"`
		Fingerprint string `json:"fingerprint"`
	}

	CreateSessionRs struct {
		AccessToken string `json:"access_token"`
	}

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
	r.HandleFunc(login, h.login).Methods(http.MethodPost)
	r.HandleFunc(refresh, h.refresh).Methods(http.MethodPost)
	r.HandleFunc(logout, middleware.Auth(h.logout)).Methods(http.MethodPost)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	handler := handler.NewHandler(w, r)
	authorization, err := handler.GetHeaderAuthorization()
	if err != nil {
		handler.SendError(http.StatusBadRequest, fmt.Errorf("auth.delivery.Handler.login: %v", err))
		return
	}
	var data CreateSessionRq
	err = decodeBasicAuth(authorization, &data)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.login - check body: %v", err))
		return
	}
	data.Fingerprint, err = handler.GetHeaderFingerprint()
	if err != nil {
		handler.SendError(http.StatusBadRequest, fmt.Errorf("auth.delivery.Handler. - GetHeaderFingerprint: %v", err))
		return
	}
	ctx := r.Context()
	tokens, err := h.authSvc.Login(ctx, data.User, data.Password, data.Fingerprint)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.login - create session: %v", err))
		return
	}

	b, err := json.Marshal(&CreateSessionRs{
		AccessToken: tokens.AccessToken,
	})
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.login - marshal: %v", err))
		return
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second
	handler.SetHeader("Content-Type", "application/json")
	handler.SetCookieRefreshToken(tokens.RefreshToken, duration)
	handler.SendData(b)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	handler := handler.NewHandler(w, r)
	refreshToken, err := handler.Cookie(common.RefreshToken)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - get cookie: %v", err))
		return
	}

	refreshID, err := uuid.Parse(refreshToken.Value)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - get cookie: %v", err))
		return
	}

	fingerprint, err := handler.GetHeaderFingerprint()
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - get fingerprint: %v", err))
		return
	}

	ctx := r.Context()

	tokens, err := h.authSvc.RefreshSessionToken(ctx, refreshID, fingerprint)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - RefreshSessionToken: %v", err))
		return
	}
	b, err := json.Marshal(&CreateSessionRs{
		AccessToken: tokens.AccessToken,
	})
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - marshal: %v", err))
		return
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second
	handler.SetHeader("Content-Type", "application/json")
	handler.SetCookieRefreshToken(tokens.RefreshToken, duration)
	handler.SendData(b)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	handler := handler.NewHandler(w, r)
	ctx := r.Context()
	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		handler.SendError(http.StatusUnauthorized, fmt.Errorf("auth.delivery.Handler.logout - unauthorized: %v", err))
		return
	}

	err = h.authSvc.Logout(ctx, uid)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.logout - logout: %v", err))
		return
	}
	handler.SendData([]byte("done"))
}

func decodeBasicAuth(auth string, data *CreateSessionRq) error {
	base, err := base64.StdEncoding.DecodeString(strings.Split(auth, " ")[1])
	if err != nil {
		return fmt.Errorf("auth.delivery.decodeBasicAuth - decode base64: %v", err)
	}
	authData := strings.Split(string(base), ":")
	if len(authData) != 2 {
		return fmt.Errorf("auth.delivery.decodeBasicAuth - invalid auth data")
	}

	data.User = authData[0]
	data.Password = authData[1]

	return nil
}

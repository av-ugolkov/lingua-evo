package delivery

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"lingua-evo/internal/config"
	"lingua-evo/internal/services/auth/dto"
	"lingua-evo/internal/services/auth/service"

	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/http/handler/header"
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

	header := header.NewHeader(w, nil)
	authorization, err := header.GetHeaderAuthorization()
	if err != nil {
		handler.SendError(w, http.StatusBadRequest, fmt.Errorf("auth.delivery.Handler.login: %v", err))
		return
	}
	var data dto.CreateSessionRq
	err = decodeBasicAuth(authorization, &data)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.login - check body: %v", err))
		return
	}
	data.Fingerprint, err = header.GetHeaderFingerprint()
	if err != nil {
		handler.SendError(w, http.StatusBadRequest, fmt.Errorf("auth.delivery.Handler.login: %v", err))
		return
	}
	ctx := r.Context()
	tokens, err := h.authSvc.Login(ctx, &data)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.login - create session: %v", err))
		return
	}

	b, err := json.Marshal(&dto.CreateSessionRs{
		AccessToken: tokens.AccessToken,
	})
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.login - marshal: %v", err))
		return
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second
	header.SetHeader("Content-Type", "application/json")
	header.SetCookieRefreshToken(tokens.RefreshToken, duration)

	_, _ = w.Write(b)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	header := header.NewHeader(w, r)
	refreshToken, err := header.GetCookieRefreshToken()
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - get cookie: %v", err))
		return
	}

	refreshID, err := uuid.Parse(refreshToken.Value)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - get cookie: %v", err))
		return
	}

	fingerprint, err := header.GetHeaderFingerprint()
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - get fingerprint: %v", err))
		return
	}

	ctx := r.Context()

	tokens, err := h.authSvc.RefreshSessionToken(ctx, refreshID, fingerprint)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - RefreshSessionToken: %v", err))
		return
	}
	b, err := json.Marshal(&dto.CreateSessionRs{
		AccessToken: tokens.AccessToken,
	})
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - marshal: %v", err))
		return
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second
	header.SetHeader("Content-Type", "application/json")
	header.SetCookieRefreshToken(tokens.RefreshToken, duration)

	_, _ = w.Write(b)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	ctx := r.Context()
	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		handler.SendError(w, http.StatusUnauthorized, fmt.Errorf("auth.delivery.Handler.logout - unauthorized: %v", err))
		return
	}

	err = h.authSvc.Logout(ctx, uid)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.logout - logout: %v", err))
		return
	}

	_, _ = w.Write([]byte("done"))
}

func decodeBasicAuth(auth string, data *dto.CreateSessionRq) error {
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
